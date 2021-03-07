package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jponc/rank-app/internal/repository"
	"github.com/jponc/rank-app/internal/s3repository"
	"github.com/jponc/rank-app/internal/uploader"
	"github.com/jponc/rank-app/pkg/dynamodb"
	"github.com/jponc/rank-app/pkg/s3manager"
	log "github.com/sirupsen/logrus"
)

func main() {
	config, err := NewConfig()
	if err != nil {
		log.Fatalf("cannot initialise config %v", err)
	}

	dynamodbClient, err := dynamodb.NewClient(config.AWSRegion, config.DBTableName)
	if err != nil {
		log.Fatalf("cannot initialise dynamodb client %v", err)
	}

	repository, err := repository.NewClient(dynamodbClient)
	if err != nil {
		log.Fatalf("cannot initialise repository %v", err)
	}

	s3ManagerClient, err := s3manager.NewClient(config.AWSRegion)
	if err != nil {
		log.Fatalf("cannot initialise s3manager client %v", err)
	}

	s3Repository, err := s3repository.NewClient(s3ManagerClient, config.S3BucketName)
	if err != nil {
		log.Fatalf("cannot initialise s3repository %v", err)
	}

	service := uploader.NewService(s3Repository, repository)

	lambda.Start(service.UploadLatestCSV)
}
