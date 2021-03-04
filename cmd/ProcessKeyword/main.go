package main

import (
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jponc/rank-app/internal/processor"
	"github.com/jponc/rank-app/internal/repository"
	"github.com/jponc/rank-app/pkg/dynamodb"
	pkgHttp "github.com/jponc/rank-app/pkg/http"
	"github.com/jponc/rank-app/pkg/zenserp"
	log "github.com/sirupsen/logrus"
)

func main() {
	config, err := NewConfig()
	if err != nil {
		log.Fatalf("cannot initialise config %v", err)
	}

	httpClient := pkgHttp.DefaultHTTPClient(time.Duration(1 * time.Minute))
	zenserpClient, err := zenserp.NewClient(config.ZenserpApiKey, zenserp.WithHTTPClient(httpClient))
	if err != nil {
		log.Fatalf("cannot initialise zenserp client %v", err)
	}

	dynamodbClient, err := dynamodb.NewClient(config.AWSRegion, config.DBTableName)
	if err != nil {
		log.Fatalf("cannot initialise dynamodb client %v", err)
	}

	repository, err := repository.NewClient(dynamodbClient)
	if err != nil {
		log.Fatalf("cannot initialise repository %v", err)
	}

	service := processor.NewService(zenserpClient, repository)

	lambda.Start(service.ProcessKeyword)
}
