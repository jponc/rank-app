package processor

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jponc/rank-app/api"
	"github.com/jponc/rank-app/internal/repository/ddbrepository"
	"github.com/jponc/rank-app/internal/repository/esrepository"
	"github.com/jponc/rank-app/internal/types"
	"github.com/jponc/rank-app/pkg/sns"
	"github.com/jponc/rank-app/pkg/zenserp"
	log "github.com/sirupsen/logrus"
)

// Service interface implements functions available for this service
type Service interface {
	// ProcessKeyword processes the keyword by sending a zenserp client search request and storing the data
	ProcessKeyword(ctx context.Context, snsEvent events.SNSEvent)
	// AddResultItemToES adds result item to elasticsearch
	AddResultItemToES(ctx context.Context, snsEvent events.SNSEvent)
}

type service struct {
	zenserpClient zenserp.Client
	ddbrepository ddbrepository.Repository
	snsClient     sns.Client
	esrepository  esrepository.Repository
}

// NewService instantiates a new service
func NewService(zenserpClient zenserp.Client, ddbrepository ddbrepository.Repository, snsClient sns.Client, esrepository esrepository.Repository) Service {
	return &service{
		zenserpClient: zenserpClient,
		ddbrepository: ddbrepository,
		snsClient:     snsClient,
		esrepository:  esrepository,
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

	resultItems := &types.ResultItemArray{}
	if err = resultItems.Unmarshal(res); err != nil {
		log.Fatalf("unable to unmarshal crawl result to result items: %v", err)
	}

	errorMsgs := []string{}

	// Iterate all result items
	for _, item := range *resultItems {
		// Store ResultItem to DB
		if err = s.ddbrepository.CreateResultItem(ctx, item); err != nil {
			errorMsgs = append(errorMsgs, err.Error())
			continue
		}

		// Send ResultItemCreated message
		msg := api.ResultItemCreatedMessage{
			ResultItem: item,
		}

		if err = s.snsClient.Publish(ctx, api.ResultItemCreated, msg); err != nil {
			errorMsgs = append(errorMsgs, err.Error())
		}
	}

	if len(errorMsgs) > 0 {
		log.Fatalf("errors: %s", strings.Join(errorMsgs, "; "))
	}

	log.Infof("crawl result successfully created for keyword: %s", processKeywordMsg.Keyword)
}

func (s *service) AddResultItemToES(ctx context.Context, snsEvent events.SNSEvent) {
	snsMsg := snsEvent.Records[0].SNS.Message

	var resultItemCreatedMsg api.ResultItemCreatedMessage
	err := json.Unmarshal([]byte(snsMsg), &resultItemCreatedMsg)
	if err != nil {
		log.Fatalf("unable to unarmarshal message: %v", err)
	}

	if s.esrepository == nil {
		log.Fatalf("esrepository not defined")
	}

	resultItem := resultItemCreatedMsg.ResultItem
	err = s.esrepository.IndexResultItem(ctx, resultItem)
	if err != nil {
		log.Fatalf("unable to index result item: %v", err)
	}
}
