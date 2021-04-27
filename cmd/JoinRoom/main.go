package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jponc/estimatex-serverless/internal/api"
	"github.com/jponc/estimatex-serverless/internal/auth"
	"github.com/jponc/estimatex-serverless/internal/repository/ddbrepository"
	"github.com/jponc/estimatex-serverless/pkg/dynamodb"
	"github.com/jponc/estimatex-serverless/pkg/sns"
)

func main() {
	config, err := NewConfig()
	if err != nil {
		log.Fatalf("cannot initialise config %v", err)
	}

	authClient, err := auth.NewClient(config.JWTSecret)
	if err != nil {
		log.Fatalf("cannot initialise auth client %v", err)
	}

	dynamodbClient, err := dynamodb.NewClient(config.AWSRegion, config.DBTableName)
	if err != nil {
		log.Fatalf("cannot initialise dynamodb client %v", err)
	}

	ddbrepository, err := ddbrepository.NewClient(dynamodbClient)
	if err != nil {
		log.Fatalf("cannot initialise ddbrepository %v", err)
	}

	snsClient, err := sns.NewClient(config.AWSRegion, config.SNSPrefix)
	if err != nil {
		log.Fatalf("cannot initialise sns client %v", err)
	}

	service := api.NewService(ddbrepository, snsClient, authClient)
	lambda.Start(service.JoinRoom)
}
