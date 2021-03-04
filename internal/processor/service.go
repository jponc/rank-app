package processor

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jponc/rank-app/api"
	"github.com/jponc/rank-app/internal/repository"
	"github.com/jponc/rank-app/pkg/zenserp"
	log "github.com/sirupsen/logrus"
)

// Service interface implements functions available for this service
type Service interface {
	ProcessKeyword(ctx context.Context, snsEvent events.SNSEvent)
}

type service struct {
	zenserpClient zenserp.Client
	repository    repository.Repository
}

// NewService instantiates a new service
func NewService(zenserpClient zenserp.Client, repository repository.Repository) Service {
	return &service{
		zenserpClient: zenserpClient,
		repository:    repository,
	}
}

// ProcessKeyword processes one keyword
func (s *service) ProcessKeyword(ctx context.Context, snsEvent events.SNSEvent) {
	snsMsg := snsEvent.Records[0].SNS.Message

	var processKeywordMsg api.ProcessKeywordMessage

	err := json.Unmarshal([]byte(snsMsg), &processKeywordMsg)
	if err != nil {
		log.Fatalf("unable to unarmarshal message: %v", err)
	}

	res, err := s.zenserpClient.Search(
		ctx,
		processKeywordMsg.Keyword,
		processKeywordMsg.SearchEngine,
		processKeywordMsg.Device,
		processKeywordMsg.Count,
	)

	if err != nil {
		log.Fatalf("Unable to query data from zenserp using keyword: %s", processKeywordMsg.Keyword)
	}

	_, err = s.repository.CreateCrawlResult(res)
	if err != nil {
		log.Fatalf("Unable to create crawl result: %v", err)
	}

	log.Infof("Crawl result successfully created for keyword: %s", processKeywordMsg.Keyword)
}
