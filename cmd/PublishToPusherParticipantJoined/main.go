package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jponc/estimatex-serverless/internal/webhooks"
	"github.com/jponc/estimatex-serverless/pkg/pusher"
)

func main() {
	config, err := NewConfig()
	if err != nil {
		log.Fatalf("cannot initialise config %v", err)
	}

	pusherClient, err := pusher.NewClient(config.PusherAppID, config.PusherKey, config.PusherSecret, config.PusherCluster)
	if err != nil {
		log.Fatalf("cannot initialise sns client %v", err)
	}

	service := webhooks.NewService(pusherClient)
	lambda.Start(service.PublishToPusherParticipantJoined)
}
