package hello

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jponc/rank-app/pkg/lambdaresponses"
)

// Service interface implements functions available for this service
type Service interface {
	SayHello(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
}

type service struct {
	responses lambdaresponses.Responses
}

// NewService instantiates a new service
func NewService(responses lambdaresponses.Responses) Service {
	return &service{
		responses: responses,
	}
}
