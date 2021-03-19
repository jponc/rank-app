package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jponc/rank-app/internal/repository/ddbrepository"
	"github.com/jponc/rank-app/internal/repository/s3repository"
	"github.com/jponc/rank-app/internal/uploader"
	"github.com/jponc/rank-app/pkg/dynamodb"
	"github.com/jponc/rank-app/pkg/s3"
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

	ddbrepository, err := ddbrepository.NewClient(dynamodbClient)
	if err != nil {
		log.Fatalf("cannot initialise ddbrepository %v", err)
	}

	s3Client, err := s3.NewClient(config.AWSRegion)
	if err != nil {
		log.Fatalf("cannot initialise s3 client %v", err)
	}

	s3Repository, err := s3repository.NewClient(s3Client, config.S3BucketName)
	if err != nil {
		log.Fatalf("cannot initialise s3repository %v", err)
	}

	service := uploader.NewService(s3Repository, ddbrepository)

	lambda.Start(service.UploadLatestCSV)
}
