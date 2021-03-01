package hello

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jponc/rank-app/api"
)

// SaysHello says hello
func (s *service) SayHello(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
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