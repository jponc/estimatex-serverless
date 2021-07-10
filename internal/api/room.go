package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jponc/estimatex-serverless/api/schema"
	"github.com/jponc/estimatex-serverless/internal/repository/ddbrepository"
	"github.com/jponc/estimatex-serverless/pkg/lambdaresponses"
	log "github.com/sirupsen/logrus"
)

func (s *Service) HostRoom(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if s.authClient == nil || s.ddbrepository == nil {
		log.Errorf("authClient or ddbrepository is nil")
		return lambdaresponses.Respond500()
	}

	req := &schema.HostRoomRequest{}

	err := json.Unmarshal([]byte(request.Body), req)
	if err != nil {
		return lambdaresponses.Respond400(fmt.Errorf("failed to unmarshal body"))
	}

	if req.Name == "" {
		return lambdaresponses.Respond400(fmt.Errorf("name can't be blank"))
	}

	// TODO Wrap both in a transaction, dynamoDB now supports transactions
	room, err := s.ddbrepository.CreateRoom(ctx)
	if err != nil {
		log.Errorf("error creating room: %w", err)
		return lambdaresponses.Respond500()
	}

	participant, err := s.ddbrepository.CreateParticipant(ctx, room.ID, req.Name, true)
	if err != nil {
		log.Errorf("error creating participant: %w", err)
		return lambdaresponses.Respond500()
	}

	token, err := s.authClient.CreateAccessToken(*participant)
	if err != nil {
		log.Errorf("error creating access token: %w", err)
		return lambdaresponses.Respond500()
	}

	res := schema.HostRoomResponse{
		RoomID:      room.ID,
		AccessToken: token,
	}

	return lambdaresponses.Respond200(res)
}

func (s *Service) FindRoom(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if s.ddbrepository == nil {
		log.Errorf("ddbrepository is nil")
		return lambdaresponses.Respond500()
	}

	roomID, ok := request.RequestContext.Authorizer["RoomID"].(string)
	if !ok {
		return lambdaresponses.Respond500()
	}

	if roomID == "" {
		return lambdaresponses.Respond400(fmt.Errorf("roomID can't be blank"))
	}

	room, err := s.ddbrepository.FindRoom(ctx, roomID)
	if err != nil {
		if errors.Is(err, ddbrepository.ErrNotFound) {
			return lambdaresponses.Respond404(fmt.Errorf("room not found"))
		} else {

			log.Errorf("error finding room: %w", err)
			return lambdaresponses.Respond500()
		}
	}

	res := schema.FindRoomResponse{
		Room: *room,
	}

	return lambdaresponses.Respond200(res)
}

func (s *Service) JoinRoom(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if s.ddbrepository == nil || s.snsClient == nil || s.authClient == nil {

		log.Errorf("ddbrepository or snsClient or authClient is nil")
		return lambdaresponses.Respond500()
	}

	req := &schema.JoinRoomRequest{}

	err := json.Unmarshal([]byte(request.Body), req)
	if err != nil {
		return lambdaresponses.Respond400(fmt.Errorf("failed to unmarshal body"))
	}

	if req.RoomID == "" {
		return lambdaresponses.Respond400(fmt.Errorf("roomID can't be blank"))
	}

	if req.Name == "" {
		return lambdaresponses.Respond400(fmt.Errorf("name can't be blank"))
	}

	_, err = s.ddbrepository.FindRoom(ctx, req.RoomID)
	if err != nil {
		return lambdaresponses.Respond404(fmt.Errorf("room ID not found"))
	}

	existingParticipant, err := s.ddbrepository.FindParticipant(ctx, req.RoomID, req.Name)
	if err != nil && !errors.Is(err, ddbrepository.ErrNotFound) {

		log.Errorf("error finding participant: %w", err)
		return lambdaresponses.Respond500()
	}
	if existingParticipant != nil {
		return lambdaresponses.Respond404(fmt.Errorf("participant already exists"))
	}

	participant, err := s.ddbrepository.CreateParticipant(ctx, req.RoomID, req.Name, false)
	if err != nil {
		log.Errorf("error creating participant: %w", err)
		return lambdaresponses.Respond500()
	}

	token, err := s.authClient.CreateAccessToken(*participant)
	if err != nil {

		log.Errorf("error creating access token: %w", err)
		return lambdaresponses.Respond500()
	}

	msg := schema.ParticipantJoinedMessage{
		RoomID:          req.RoomID,
		ParticipantName: req.Name,
	}

	err = s.snsClient.Publish(ctx, schema.ParticipantJoined, msg)
	if err != nil {
		log.Errorf("error publishing participant joined to sns: %w", err)
		return lambdaresponses.Respond500()
	}

	res := schema.JoinRoomResponse{
		AccessToken: token,
	}

	return lambdaresponses.Respond200(res)
}
