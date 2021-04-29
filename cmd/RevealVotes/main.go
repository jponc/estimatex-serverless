package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jponc/estimatex-serverless/internal/api"
	"github.com/jponc/estimatex-serverless/pkg/sns"
)

func main() {
	config, err := NewConfig()
	if err != nil {
		log.Fatalf("cannot initialise config %v", err)
	}

	snsClient, err := sns.NewClient(config.AWSRegion, config.SNSPrefix)
	if err != nil {
		log.Fatalf("cannot initialise sns client %v", err)
	}

	service := api.NewService(nil, snsClient, nil)
	lambda.Start(service.RevealVotes)
}
