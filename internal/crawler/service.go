package crawler

import (
	"context"
	"fmt"

	"github.com/jponc/rank-app/api"
	"github.com/jponc/rank-app/pkg/sns"
	log "github.com/sirupsen/logrus"
)

// Service interface implements functions available for this service
type Service interface {
	RunCrawl(ctx context.Context)
}

type service struct {
	snsClient sns.Client
	keywords  *[]string
}

// NewService instantiates a new service
func NewService(snsClient sns.Client, keywords *[]string) Service {
	return &service{
		snsClient: snsClient,
		keywords:  keywords,
	}
}

// RunCrawl runs the crawler
func (s *service) RunCrawl(ctx context.Context) {
	log.Info("Crawl running")

	if s.keywords == nil {
		log.Fatalf("keywords cannot be empty")
	}

	var allErr error

	// TODO Convert ot use goroutines, waitgroups, and channels
	// TODO Make device and search engine dynamic
	for _, keyword := range *s.keywords {
		msg := api.ProcessKeywordMessage{
			Keyword:      keyword,
			Device:       "desktop",
			SearchEngine: "google.com",
			Count:        100,
		}
		err := s.snsClient.Publish(ctx, api.ProcessKeyword, msg)

		if err != nil {
			allErr = fmt.Errorf("%w; %v", allErr, err)
		}
	}

	if allErr != nil {
		log.Fatalf(allErr.Error())
	}
}
