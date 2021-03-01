package rankings

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jponc/rank-app/api"
	log "github.com/sirupsen/logrus"
)

// RunCrawl runs the crawler
func (s *service) RunCrawl(ctx context.Context, snsEvent events.SNSEvent) {
	log.Info("Crawl running")

	keywords := []string{
		"craft beers",
		"beers",
		"gin",
		"whisky",
		"rum",
	}

	var allErr error

	// TODO Convert ot use goroutines, waitgroups, and channels
	for _, keyword := range keywords {
		msg := api.ProcessKeywordMessage{Keyword: keyword, Count: 100}
		err := s.snsClient.Publish("ProcessKeyword", msg)

		if err != nil {
			allErr = fmt.Errorf("%w; Second error", err)
		}
	}

	if allErr != nil {
		log.Fatalf(allErr.Error())
	}
}
