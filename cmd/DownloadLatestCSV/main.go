package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jponc/rank-app/internal/apiservice"
	"github.com/jponc/rank-app/internal/s3repository"
	"github.com/jponc/rank-app/pkg/lambdaresponses"
	"github.com/jponc/rank-app/pkg/s3"
	log "github.com/sirupsen/logrus"
)

func main() {
	config, err := NewConfig()
	if err != nil {
		log.Fatalf("cannot initialise config %v", err)
	}

	s3Client, err := s3.NewClient(config.AWSRegion)
	if err != nil {
		log.Fatalf("cannot initialise s3 client %v", err)
	}

	s3Repository, err := s3repository.NewClient(s3Client, config.S3BucketName)
	if err != nil {
		log.Fatalf("cannot initialise s3repository %v", err)
	}

	responses := lambdaresponses.NewResponses()
	service := apiservice.NewService(responses, s3Repository)
	lambda.Start(service.DownloadLatestCSV)
}
