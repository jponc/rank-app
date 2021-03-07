package s3repository

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	awsS3Manager "github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/jponc/rank-app/pkg/s3manager"
)

const LatestCSVKey = "csv/latest.csv"

// Repository interface
type Repository interface {
	// UploadLatestCSV uploads latest csv to s3
	UploadLatestCSV(csvData string) error
}

type repository struct {
	s3Manager  s3manager.Client
	bucketName string
}

// NewClient instantiates a repository
func NewClient(s3Manager s3manager.Client, bucketName string) (Repository, error) {
	r := &repository{
		s3Manager:  s3Manager,
		bucketName: bucketName,
	}

	return r, nil
}

func (r *repository) UploadLatestCSV(csvData string) error {
	input := &awsS3Manager.UploadInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(LatestCSVKey),
		Body:   strings.NewReader(csvData),
	}

	_, err := r.s3Manager.Upload(input)
	if err != nil {
		return fmt.Errorf("failed to upload latest csv to s3: %v", err)
	}

	return nil
}
