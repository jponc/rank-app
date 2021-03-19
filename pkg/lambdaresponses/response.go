package lambdaresponses

import (
	"github.com/aws/aws-lambda-go/events"
	"gopkg.in/square/go-jose.v2/json"
)

type errorResponseBody struct {
	Error string `json:"error"`
}

func Respond500() (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{Body: "Internal Server Error", StatusCode: 500}, nil
}

func Respond400(err error) (events.APIGatewayProxyResponse, error) {
	resBody := errorResponseBody{
		Error: err.Error(),
	}

	body, err := json.Marshal(resBody)
	if err != nil {
		return Respond500()
	}

	return events.APIGatewayProxyResponse{Body: string(body), StatusCode: 400}, nil
}

func Respond200(body interface{}) (events.APIGatewayProxyResponse, error) {
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return Respond500()
	}

	return events.APIGatewayProxyResponse{Body: string(bodyJson), StatusCode: 200}, nil
}

func Respond302(location string) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: 302,
		Headers: map[string]string{
			"Location": location,
		},
	}, nil
}
