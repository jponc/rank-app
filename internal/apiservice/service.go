package apiservice

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jponc/rank-app/internal/s3repository"
	"github.com/jponc/rank-app/pkg/lambdaresponses"
)

// Service interface implements functions available for this service
type Service interface {
	DownloadLatestCSV(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
}

type service struct {
	responses    lambdaresponses.Responses
	s3repository s3repository.Repository
}

// NewService instantiates a new service
func NewService(responses lambdaresponses.Responses, s3repository s3repository.Repository) Service {
	return &service{
		responses:    responses,
		s3repository: s3repository,
	}
}

func (s *service) DownloadLatestCSV(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if s.s3repository == nil {
		return s.responses.Respond500()
	}

	urlStr, err := s.s3repository.GetURLLatestCSV()
	if err != nil {
		return s.responses.Respond500()
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 302,
		Headers: map[string]string{
			"Location": urlStr,
		},
	}, nil
}
