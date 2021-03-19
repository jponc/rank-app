package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jponc/rank-app/internal/repository/esrepository"
	"github.com/jponc/rank-app/pkg/elasticsearch"
	log "github.com/sirupsen/logrus"
)

func main() {
	config, err := NewConfig()
	if err != nil {
		log.Fatalf("cannot initialise config %v", err)
	}

	esClient, err := elasticsearch.NewClient(config.ElasticsearchURL, config.AWSRegion)
	if err != nil {
		log.Fatalf("cannot initialise elasticsearch client %v", err)
	}

	esRepository, err := esrepository.NewRepository(esClient)
	if err != nil {
		log.Fatalf("cannot initialise es repository %v", err)
	}

	lambda.Start(esRepository.UpdateResultItemIndexMapping)
}
