package apiservice

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jponc/rank-app/internal/repository"
	"github.com/jponc/rank-app/pkg/lambdaresponses"
)

// Service interface implements functions available for this service
type Service interface {
	SayHello(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
}

type service struct {
	responses  lambdaresponses.Responses
	repository repository.Repository
}

// NewService instantiates a new service
func NewService(responses lambdaresponses.Responses, repository repository.Repository) Service {
	return &service{
		responses:  responses,
		repository: repository,
	}
}

func (s *service) SayHello(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if s.repository == nil {
		return s.responses.Respond500()
	}

	crawlResults, err := s.repository.GetLatestCrawlResults()
	if err != nil {
		return s.responses.Respond500()
	}

	pairs := [][]string{
		{
			"query",
			"search_engine",
			"device",
			"item_position",
			"item_title",
			"item_url",
			"item_description",
			"result_url",
			"result_id",
			"created_at",
		},
	}

	for _, crawlResult := range *crawlResults {
		for _, item := range crawlResult.Items {
			row := []string{
				crawlResult.Query,
				crawlResult.SearchEngine,
				crawlResult.Device,
				fmt.Sprint(item.Position),
				item.Title,
				item.URL,
				item.Description,
				crawlResult.URL,
				crawlResult.ID.String(),
				crawlResult.CreatedAt.String(),
			}

			pairs = append(pairs, row)
		}
	}

	b := new(bytes.Buffer)
	w := csv.NewWriter(b)
	w.WriteAll(pairs)

	if w.Error() != nil {
		return s.responses.Respond500()
	}

	csvString := b.String()
	now := time.Now()
	filename := fmt.Sprintf("crawl_results_%d-%d-%d_.csv", now.Year(), now.Month(), now.Day())

	return events.APIGatewayProxyResponse{
		Body:       csvString,
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type":        "text/csv",
			"Content-disposition": fmt.Sprintf("attachment; filename=%s", filename),
		},
	}, nil
}
