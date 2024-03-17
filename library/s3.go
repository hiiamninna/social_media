package library

import (
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3Config struct {
	ID         string
	SecretKey  string
	BucketName string
}

type S3 struct {
	Uploader *s3manager.Uploader
	Config   S3Config
}

func NewS3(config S3Config) (S3, error) {
	session, err := session.NewSession(&aws.Config{
		Region:      aws.String("ap-southeast-1"),
		Credentials: credentials.NewStaticCredentials(config.ID, config.SecretKey, ""),
	})
	if err != nil {
		return S3{}, fmt.Errorf("open con : %w", err)
	}

	return S3{
		Uploader: s3manager.NewUploader(session),
		Config:   config,
	}, nil
}

func (storage *S3) UploadFile(data io.Reader, FileName string) (string, error) {
	resp, err := storage.Uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(storage.Config.BucketName),
		Key:    aws.String(FileName),
		Body:   data,
		ACL:    aws.String("public-read"),
	})
	if err != nil {
		return "", fmt.Errorf("fail upload file : %w", err)
	}

	return resp.Location, nil
}
