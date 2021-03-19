package s3repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	awsS3 "github.com/aws/aws-sdk-go/service/s3"
	awsS3Manager "github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/jponc/rank-app/pkg/s3"
)

const latestCSVKey = "csv/latest.csv"

// Repository interface
type Repository interface {
	// UploadLatestCSV uploads latest csv to s3
	UploadLatestCSV(csvData string) error
	// GetURLLatestCSV gets the url of latest csv
	GetURLLatestCSV() (string, error)
}

type repository struct {
	s3Client   s3.Client
	bucketName string
}

// NewClient instantiates a repository
func NewClient(s3Client s3.Client, bucketName string) (Repository, error) {
	r := &repository{
		s3Client:   s3Client,
		bucketName: bucketName,
	}

	return r, nil
}

func (r *repository) UploadLatestCSV(csvData string) error {
	input := &awsS3Manager.UploadInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(latestCSVKey),
		Body:   strings.NewReader(csvData),
	}

	_, err := r.s3Client.Upload(input)
	if err != nil {
		return fmt.Errorf("failed to upload latest csv to s3: %v", err)
	}

	return nil
}

func (r *repository) GetURLLatestCSV() (string, error) {
	input := &awsS3.GetObjectInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(latestCSVKey),
	}

	req, _ := r.s3Client.GetObjectRequest(input)
	urlStr, err := req.Presign(15 * time.Minute)

	if err != nil {
		return "", fmt.Errorf("failed to get latest csv url: %v", err)
	}

	return urlStr, nil
}
