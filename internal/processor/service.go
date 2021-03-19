package processor

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jponc/rank-app/api"
	"github.com/jponc/rank-app/internal/repository/ddbrepository"
	"github.com/jponc/rank-app/pkg/sns"
	"github.com/jponc/rank-app/pkg/zenserp"
	log "github.com/sirupsen/logrus"
)

// Service interface implements functions available for this service
type Service interface {
	// ProcessKeyword processes the keyword by sending a zenserp client search request and storing the data
	ProcessKeyword(ctx context.Context, snsEvent events.SNSEvent)
	// AddCrawlResultToLatest adds CrawlResult to latest PK
	AddCrawlResultToLatest(ctx context.Context, snsEvent events.SNSEvent)
}

type service struct {
	zenserpClient zenserp.Client
	ddbrepository ddbrepository.Repository
	snsClient     sns.Client
}

// NewService instantiates a new service
func NewService(zenserpClient zenserp.Client, ddbrepository ddbrepository.Repository, snsClient sns.Client) Service {
	return &service{
		zenserpClient: zenserpClient,
		ddbrepository: ddbrepository,
		snsClient:     snsClient,
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

	if s.zenserpClient == nil {
		log.Fatalf("zenserpClient not defined")
	}

	if s.ddbrepository == nil {
		log.Fatalf("ddbrepository not defined")
	}

	res, err := s.zenserpClient.Search(
		ctx,
		processKeywordMsg.Keyword,
		processKeywordMsg.SearchEngine,
		processKeywordMsg.Device,
		processKeywordMsg.Count,
	)

	if err != nil {
		log.Fatalf("unable to query data from zenserp using keyword: %s", processKeywordMsg.Keyword)
	}

	crawlResult, err := s.ddbrepository.CreateCrawlResult(res)
	if err != nil {
		log.Fatalf("unable to create crawl result: %v", err)
	}

	msg := api.CrawlResultCreatedMessage{
		CrawlResult: *crawlResult,
	}

	err = s.snsClient.Publish(api.CrawlResultCreated, msg)

	log.Infof("crawl result successfully created for keyword: %s", processKeywordMsg.Keyword)
}

func (s *service) AddCrawlResultToLatest(ctx context.Context, snsEvent events.SNSEvent) {
	snsMsg := snsEvent.Records[0].SNS.Message

	var crawlResultCreatedMsg api.CrawlResultCreatedMessage
	err := json.Unmarshal([]byte(snsMsg), &crawlResultCreatedMsg)
	if err != nil {
		log.Fatalf("unable to unarmarshal message: %v", err)
	}

	if s.ddbrepository == nil {
		log.Fatalf("ddbrepository not defined")
	}

	crawlResult := crawlResultCreatedMsg.CrawlResult

	err = s.ddbrepository.AddCrawlResultToLatest(&crawlResult)
	if err != nil {
		log.Fatalf("unable to store crawl result to latest: %v", err)
	}

	log.Infof("latest data for %s %s %s updated", crawlResult.Query, crawlResult.SearchEngine, crawlResult.Device)
}
