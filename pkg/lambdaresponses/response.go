package lambdaresponses

import (
	"github.com/aws/aws-lambda-go/events"
	"gopkg.in/square/go-jose.v2/json"
)

type errorResponseBody struct {
	Error string `json:"error"`
}

type Responses interface {
	// Respond500 responds 500 internal server error
	Respond500() (events.APIGatewayProxyResponse, error)
	// Respond400 responds 400 bad request
	Respond400(err error) (events.APIGatewayProxyResponse, error)
	// Respond200 responses 200 success
	Respond200(body interface{}) (events.APIGatewayProxyResponse, error)
}

type responses struct{}

// NewResponses instantiates a response
func NewResponses() Responses {
	return &responses{}
}

func (r *responses) Respond500() (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{Body: "Internal Server Error", StatusCode: 500}, nil
}

func (r *responses) Respond400(err error) (events.APIGatewayProxyResponse, error) {
	resBody := errorResponseBody{
		Error: err.Error(),
	}

	body, err := json.Marshal(resBody)
	if err != nil {
		return r.Respond500()
	}

	return events.APIGatewayProxyResponse{Body: string(body), StatusCode: 400}, nil
}

func (r *responses) Respond200(body interface{}) (events.APIGatewayProxyResponse, error) {
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return r.Respond500()
	}

	return events.APIGatewayProxyResponse{Body: string(bodyJson), StatusCode: 200}, nil
}
