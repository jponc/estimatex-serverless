package api

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jponc/estimatex-serverless/pkg/lambdaresponses"
	log "github.com/sirupsen/logrus"
)

func (s *Service) FindParticipants(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if s.ddbrepository == nil {
		log.Errorf("ddbrepository is nil")
		return lambdaresponses.Respond500()
	}

	roomID, ok := request.RequestContext.Authorizer["RoomID"].(string)
	if !ok {
		return lambdaresponses.Respond500()
	}

	participants, err := s.ddbrepository.FindParticipants(ctx, roomID)
	if err != nil {
		log.Errorf("error finding participants: %w", err)
		return lambdaresponses.Respond500()
	}

	return lambdaresponses.Respond200(participants)
}
