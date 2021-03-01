package rankings

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jponc/rank-app/pkg/lambdaresponses"
	"github.com/jponc/rank-app/pkg/sns"
	"github.com/jponc/rank-app/pkg/zenserp"
)

// Service interface implements functions available for this service
type Service interface {
	RunCrawl(ctx context.Context, snsEvent events.SNSEvent)
	ProcessKeyword(ctx context.Context, snsEvent events.SNSEvent)
}

type service struct {
	responses     lambdaresponses.Responses
	zenserpClient zenserp.Client
	snsClient     sns.Client
}

// NewService instantiates a new service
func NewService(responses lambdaresponses.Responses, zenserpClient zenserp.Client, snsClient sns.Client) Service {
	return &service{
		responses:     responses,
		zenserpClient: zenserpClient,
		snsClient:     snsClient,
	}
}
