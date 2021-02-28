package rankings

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jponc/rank-app/api"
	"github.com/jponc/rank-app/pkg/lambdaresponses"
)

// Service interface implements functions available for this service
type Service interface {
	SayHello(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
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

// SaysHello says hello
func (s *service) SayHello(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	req := &api.SayHelloRequest{}

	err := json.Unmarshal([]byte(request.Body), req)
	if err != nil {
		return s.responses.Respond400(fmt.Errorf("Failed to unmarshall body"))
	}

	if req.Name == "Waldo" {
		return s.responses.Respond400(fmt.Errorf("Cannot use name Waldo!"))
	}

	message := fmt.Sprintf("Hello %s", req.Name)
	return s.responses.Respond200(api.SayHelloResponse{Message: message})
}
