package syncer

import (
	"backup-to-s3/fileLister"
	"backup-to-s3/logging"
	"backup-to-s3/options"
	"backup-to-s3/s3client"
	"backup-to-s3/supervisor"
	"backup-to-s3/timing"
	"fmt"
	"os"
	"os/signal"
	"sort"
	"syscall"
	"time"
)

type Syncer struct {
	lastFiles  map[string]bool
	logger     logging.Logger
	fileLister fileLister.FileLister
	s3client   s3client.S3Client
	opts       *options.Options
	supervisor supervisor.Supervisor
	timer      timing.ExponentialTimer
}

var optionsPath = "/data/options.json"

func New(logger logging.Logger) *Syncer {
	s := &Syncer{
		lastFiles: map[string]bool{},
		logger:    logger,
	}

	opts, err := options.New(optionsPath)
	if err != nil {
		s.halt(err)
	}

	s.opts = opts

	s.fileLister = fileLister.New(s.opts)
	s.supervisor = supervisor.New(s.opts)
	s.timer = timing.New(s.opts.FilesCheckInterval, 10*time.Minute)

	s.s3client, err = s3client.New(s.opts, s.logger)
	if err != nil {
		s.halt(err)
	}

	return s
}

func (s *Syncer) Run() {
	s.logger.Info("Addon ready")

	for {
		currFiles, err := s.fileLister.ListCurrFiles()
		if err != nil {
			s.logger.Error(err)
			s.timer.Failed()
			continue
		}

		if s.hasNewFiles(currFiles) {
			err = s.syncBackups(currFiles)
			if err != nil {
				s.logger.Error(err)
				s.timer.Failed()
				continue
			}
		}

		s.timer.Succeeded()
	}
}

func (s *Syncer) syncBackups(currFiles []string) error {
	s3files, err := s.s3client.ListS3Files()
	if err != nil {
		return err
	}

	haBackupFiles, err := s.supervisor.ListHaApiFiles()
	if err != nil {
		return err
	}

	// Sort to have predictable processing order for tests
	sort.Strings(currFiles)

	for _, localFilename := range currFiles {
		filenameKey, _ := s.fileLister.GetNewFormat(localFilename)
		s3Filename, inDb := haBackupFiles[filenameKey]

		if !inDb {
			s.logger.Warn(fmt.Sprintf("File '%s' doesn't have corresponding record in DB yet, will try sync again later",
				localFilename))
			return nil
		}

		if !s3files[s3Filename] {
			s.logger.Info(fmt.Sprintf("Uploading '%s' as '%s'", localFilename, s3Filename))
			err = s.s3client.Upload(localFilename, s3Filename)
			if err != nil {
				return err
			}
			s.lastFiles[localFilename] = true
		}
	}

	// The addon may run for a long time. Reassign lastFiles to prevent accumulating names of local backup files that
	// may have been already deleted from disk.
	s.lastFiles = map[string]bool{}
	for _, filename := range currFiles {
		s.lastFiles[filename] = true
	}

	return nil
}

func (s *Syncer) hasNewFiles(currFiles []string) bool {
	for _, filename := range currFiles {
		if _, exists := s.lastFiles[filename]; !exists {
			return true
		}
	}
	return false
}

func (s *Syncer) halt(err error) {
	if err != nil {
		s.logger.Error(err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	s.logger.Info("Halting")
	<-done
}
