package syncer

import (
	"backup-to-s3/logging"
	"backup-to-s3/options"
	"backup-to-s3/s3client"
	"backup-to-s3/supervisor"
	"backup-to-s3/timing"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Syncer struct {
	lastFiles  map[string]bool
	logger     logging.Logger
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
		currFiles, err := s.supervisor.ListHaApiFiles()
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

		break
	}
}

func (s *Syncer) syncBackups(currFiles []*supervisor.Result) error {
	s3files, err := s.s3client.ListS3Files()
	if err != nil {
		return err
	}

	for _, currFile := range currFiles {
		if !s3files[currFile.S3Filename] {
			s.logger.Info(fmt.Sprintf("Downloading '%s'", currFile.S3Filename))
			body, err2 := s.supervisor.Download(currFile.Slug)
			if err2 != nil {
				return err2
			}
			s.logger.Info(fmt.Sprintf("Uploading '%s'", currFile.S3Filename))
			err2 = s.s3client.Upload(body, currFile.S3Filename)
			if err2 != nil {
				return err2
			}
			s.lastFiles[currFile.S3Filename] = true
		}
	}

	// The addon may run for a long time. If sync succeeds, reassign lastFiles to prevent accumulating names of HA
	// backup files that may have been already deleted from HA.
	s.lastFiles = map[string]bool{}
	for _, currFile := range currFiles {
		s.lastFiles[currFile.S3Filename] = true
	}

	return nil
}

func (s *Syncer) hasNewFiles(currFiles []*supervisor.Result) bool {
	for _, currFile := range currFiles {
		if _, exists := s.lastFiles[currFile.S3Filename]; !exists {
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
