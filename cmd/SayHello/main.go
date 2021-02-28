package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jponc/rank-app/internal/rankings"
	"github.com/jponc/rank-app/pkg/lambdaresponses"
)

func main() {
	responses := lambdaresponses.NewResponses()
	service := rankings.NewService(responses)
	lambda.Start(service.SayHello)
}
