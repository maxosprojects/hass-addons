package s3client

import (
	"backup-to-s3/logging"
	"backup-to-s3/options"
	"bytes"
	"context"
	"crypto/md5"
	"encoding/base64"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io"
	"os"
	"path/filepath"
)

type S3Client interface {
	ListS3Files() (map[string]bool, error)
	Upload(localFilename, s3Filename string) error
}

type client struct {
	opts     *options.Options
	s3Client *s3.Client
	logger   logging.Logger
}

func New(opts *options.Options, logger logging.Logger) (S3Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), func(awsOpts *config.LoadOptions) error {
		awsOpts.Region = opts.BucketRegion
		awsOpts.Credentials = credentials.NewStaticCredentialsProvider(opts.AwsAccessKey, opts.AwsSecretAccessKey, "")
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &client{
		logger:   logger,
		opts:     opts,
		s3Client: s3.NewFromConfig(cfg),
	}, nil
}

func (c *client) ListS3Files() (map[string]bool, error) {
	s3files := map[string]bool{}
	c.logger.Info("Listing S3...")

	var continueToken *string
	for {
		resp, err := c.s3Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
			Bucket:            aws.String(c.opts.BucketName),
			ContinuationToken: continueToken,
		})
		if err != nil {
			return nil, err
		}

		for _, cont := range resp.Contents {
			s3files[*cont.Key] = true
		}

		if resp.ContinuationToken == nil {
			break
		}
	}

	return s3files, nil
}

func (c *client) Upload(localFilename, s3Filename string) error {
	fullFilePath := filepath.Join(c.opts.BackupDir, localFilename)

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

	_, err = c.s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:        aws.String(c.opts.BucketName),
		Key:           aws.String(s3Filename),
		Body:          bytes.NewReader(data),
		ContentLength: aws.Int64(stat.Size()),
		ContentType:   aws.String("application/x-tar"),
		ContentMD5:    aws.String(md5SumBase64),
	})

	return err
}
