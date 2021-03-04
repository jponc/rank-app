package crawler

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jponc/rank-app/api"
	"github.com/jponc/rank-app/pkg/sns"
	log "github.com/sirupsen/logrus"
)

// Service interface implements functions available for this service
type Service interface {
	RunCrawl(ctx context.Context, snsEvent events.SNSEvent)
}

type service struct {
	snsClient sns.Client
}

// NewService instantiates a new service
func NewService(snsClient sns.Client) Service {
	return &service{
		snsClient: snsClient,
	}
}

// RunCrawl runs the crawler
func (s *service) RunCrawl(ctx context.Context, snsEvent events.SNSEvent) {
	log.Info("Crawl running")

	keywords := []string{
		"craft beers",
	}

	var allErr error

	// TODO Convert ot use goroutines, waitgroups, and channels
	for _, keyword := range keywords {
		msg := api.ProcessKeywordMessage{
			Keyword:      keyword,
			Device:       "desktop",
			SearchEngine: "google.com",
			Count:        100,
		}
		err := s.snsClient.Publish("ProcessKeyword", msg)

		if err != nil {
			allErr = fmt.Errorf("%w; Second error", err)
		}
	}

	if allErr != nil {
		log.Fatalf(allErr.Error())
	}
}
