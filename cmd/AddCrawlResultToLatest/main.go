package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jponc/rank-app/internal/processor"
	"github.com/jponc/rank-app/internal/repository"
	"github.com/jponc/rank-app/pkg/dynamodb"
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

	service := processor.NewService(nil, repository, nil)

	lambda.Start(service.AddCrawlResultToLatest)
}
