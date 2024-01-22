package main

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

type Syncer struct {
	options  *Options
	s3Client *s3.Client

	s3files                     map[string]bool
	lastFiles                   map[string]bool
	sluggedFileNameToNormalName map[string]string
}

type ResponseData struct {
	Result string `json:"result"`
	Data   struct {
		Backups []struct {
			Slug       string      `json:"slug"`
			Name       string      `json:"name"`
			Date       time.Time   `json:"date"`
			Type       string      `json:"type"`
			Size       float64     `json:"size"`
			Location   interface{} `json:"location"`
			Protected  bool        `json:"protected"`
			Compressed bool        `json:"compressed"`
			Content    struct {
				Homeassistant bool     `json:"homeassistant"`
				Addons        []string `json:"addons"`
				Folders       []string `json:"folders"`
			} `json:"content"`
		} `json:"backups"`
	} `json:"data"`
}

const SupervisorTokenVarName = "SUPERVISOR_TOKEN"

var supervisorApi = "http://supervisor/backups"
var backupDir = "/backup"

func main() {
	syncer := &Syncer{}
	syncer.run()
}

func (s *Syncer) init() error {
	options, err := NewOptions()
	s.options = options
	if err != nil {
		halt(err)
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(), func(awsOpts *config.LoadOptions) error {
		awsOpts.Region = options.BucketRegion
		awsOpts.Credentials = credentials.
			NewStaticCredentialsProvider(options.AwsAccessKey, options.AwsSecretAccessKey, "")
		return nil
	})
	if err != nil {
		halt(err)
	}
	s3Client := s3.NewFromConfig(cfg)

	s.s3Client = s3Client

	return err
}

func (s *Syncer) run() {
	err := s.init()
	if err != nil {
		halt(err)
	}

	fmt.Println("Addon ready")

	timer := NewExponentialTimer(s.options.FilesCheckInterval, 10*time.Minute)

	for {
		currFiles, err := s.listCurrFiles()
		if err != nil {
			logError(err)
			timer.Failed()
			continue
		}

		if s.hasNewFiles(currFiles) {
			err = s.syncBackups(currFiles)
			if err != nil {
				logError(err)
				timer.Failed()
				continue
			}
		}

		timer.Succeeded()
	}
}

func (s *Syncer) syncBackups(currFiles map[string]bool) error {
	err := s.listS3Files()
	if err != nil {
		return err
	}

	err = s.listHaApiFiles()
	if err != nil {
		return err
	}

	for currFile, _ := range currFiles {
		filename := s.sluggedFileNameToNormalName[currFile]
		if !s.s3files[filename] {
			log.Printf("Uploading '%s' as '%s'", currFile, filename)
			err = s.upload(currFile, filename)
			if err != nil {
				return err
			}
		}
	}

	s.lastFiles = currFiles

	return nil
}

func (s *Syncer) upload(fileSlug, filename string) error {
	fullFilePath := filepath.Join(backupDir, fileSlug)

	file, err := os.Open(fullFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return err
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	md5Sum := md5.Sum(data)
	md5SumBase64 := base64.StdEncoding.EncodeToString(md5Sum[:])

	_, err = s.s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:        aws.String(s.options.BucketName),
		Key:           aws.String(filename),
		Body:          bytes.NewReader(data),
		ContentLength: aws.Int64(stat.Size()),
		ContentType:   aws.String("application/x-tar"),
		ContentMD5:    aws.String(md5SumBase64),
	})

	return err
}

func (s *Syncer) listCurrFiles() (map[string]bool, error) {
	entries, err := os.ReadDir(backupDir)
	if err != nil {
		return nil, err
	}

	res := map[string]bool{}
	for _, entry := range entries {
		if entry.Type().IsRegular() {
			res[entry.Name()] = true
		}
	}

	return res, err
}

func (s *Syncer) listHaApiFiles() error {
	req, err := http.NewRequest("GET", supervisorApi, nil)
	if err != nil {
		return err
	}

	supervisorToken, tokenPresent := os.LookupEnv(SupervisorTokenVarName)
	if !tokenPresent {
		return errors.New("supervisor token is not present")
	}

	req.Header.Set("Authorization", "Bearer "+supervisorToken)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected HTTP code from Supervisor API: %d", resp.StatusCode)
	}

	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var responseData *ResponseData

	err = json.Unmarshal(res, &responseData)
	if err != nil {
		return err
	}

	s.sluggedFileNameToNormalName = map[string]string{}
	for _, back := range responseData.Data.Backups {
		s.sluggedFileNameToNormalName[back.Slug+".tar"] = back.Name + ".tar"
	}

	return nil
}

func (s *Syncer) listS3Files() error {
	s.s3files = map[string]bool{}

	var continueToken *string
	for {
		resp, err := s.s3Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
			Bucket:            aws.String(s.options.BucketName),
			ContinuationToken: continueToken,
		})
		if err != nil {
			return err
		}

		for _, cont := range resp.Contents {
			s.s3files[*cont.Key] = true
		}

		if resp.ContinuationToken == nil {
			break
		}
	}

	return nil
}

func (s *Syncer) hasNewFiles(currFiles map[string]bool) bool {
	for name, _ := range currFiles {
		if !s.lastFiles[name] {
			return true
		}
	}
	return false
}

func halt(err error) {
	if err != nil {
		logError(err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	logInfo("Halting")
	<-done
}

func logError(err error) {
	log.Printf("[backup-to-s3 addon] ERROR %v", err)
}

func logInfo(msg string) {
	log.Printf("[backup-to-s3 addon] INFO %s", msg)
}
