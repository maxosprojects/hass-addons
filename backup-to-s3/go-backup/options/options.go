package options

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"time"
)

type Options struct {
	FilesCheckInterval time.Duration
	AwsAccessKey       string `json:"aws_access_key"`
	AwsSecretAccessKey string `json:"aws_secret_access_key"`
	BucketName         string `json:"bucket_name"`
	BucketRegion       string `json:"bucket_region"`
	StorageClass       string `json:"storage_class"`

	BackupDir     string `json:"backup_dir"`
	SupervisorApi string `json:"supervisor_api"`

	// Temporary location for unmarshalling. Use FilesCheckInterval instead
	FilesCheckIntervalStr string `json:"files_check_interval"`
}

func New(optionsPath string) (*Options, error) {
	jsonFile, err := os.Open(optionsPath)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	options := &Options{}
	err = json.Unmarshal(byteValue, options)
	if err != nil {
		return nil, err
	}

	options.FilesCheckInterval, err = time.ParseDuration(options.FilesCheckIntervalStr)
	if err != nil {
		return nil, err
	}

	if options.FilesCheckInterval == 0 {
		return nil, errors.New("interval must not be ero")
	}

	return options, nil
}
