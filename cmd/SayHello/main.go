package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jponc/rank-app/internal/hello"
)

func main() {
	service := hello.NewService()
	lambda.Start(service.SayHello)
}
