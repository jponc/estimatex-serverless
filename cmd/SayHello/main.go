package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jponc/estimatex-serverless/internal/api"
)

func main() {
	service := api.NewService(nil, nil, nil)
	lambda.Start(service.SayHello)
}
