package main

import (
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jponc/rank-app/internal/rankings"
	pkgHttp "github.com/jponc/rank-app/pkg/http"
	"github.com/jponc/rank-app/pkg/lambdaresponses"
	"github.com/jponc/rank-app/pkg/sns"
	"github.com/jponc/rank-app/pkg/zenserp"
	log "github.com/sirupsen/logrus"
)

func main() {
	responses := lambdaresponses.NewResponses()

	config, err := NewConfig()
	if err != nil {
		log.Fatalf("cannot initialise config %v", err)
	}

	snsClient, err := sns.NewClient(config.AWSRegion, config.SNSPrefix)
	if err != nil {
		log.Fatalf("cannot initialise sns client %v", err)
	}

	httpClient := pkgHttp.DefaultHTTPClient(time.Duration(1 * time.Minute))
	zenserpClient, err := zenserp.NewClient(config.ZenserpApiKey, zenserp.WithHTTPClient(httpClient))
	if err != nil {
		log.Fatalf("cannot initialise zenserp client %v", err)
	}

	service := rankings.NewService(responses, zenserpClient, snsClient)

	lambda.Start(service.RunCrawl)
}
