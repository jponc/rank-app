package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jponc/rank-app/internal/crawler"
	"github.com/jponc/rank-app/pkg/sns"
	log "github.com/sirupsen/logrus"
)

func main() {
	config, err := NewConfig()
	if err != nil {
		log.Fatalf("cannot initialise config %v", err)
	}

	snsClient, err := sns.NewClient(config.AWSRegion, config.SNSPrefix)
	if err != nil {
		log.Fatalf("cannot initialise sns client %v", err)
	}

	service := crawler.NewService(snsClient, &config.Keywords)

	lambda.Start(service.RunCrawl)
}
