package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jponc/estimatex-serverless/api/schema"
	"github.com/jponc/estimatex-serverless/internal/auth"
	"github.com/jponc/estimatex-serverless/internal/repository/ddbrepository"
	"github.com/jponc/estimatex-serverless/pkg/lambdaresponses"
	"github.com/jponc/estimatex-serverless/pkg/sns"
	log "github.com/sirupsen/logrus"
)

// Service interface implements functions available for this service
type Service interface {
	// SayHello is a placeholder endpoint
	SayHello(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
	// HostRoom creates a new room
	HostRoom(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
	// FindRoom finds the room given a room ID
	FindRoom(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
	// JoinRoom allows users to join an existing room
	JoinRoom(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
	// CastVote allows the participant to cast a vote
	CastVote(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
}

type service struct {
	ddbrepository ddbrepository.Repository
	snsClient     sns.Client
	authClient    auth.Client
}

// NewService instantiates a new service
func NewService(ddbrepository ddbrepository.Repository, snsClient sns.Client, authClient auth.Client) Service {
	return &service{
		ddbrepository: ddbrepository,
		snsClient:     snsClient,
		authClient:    authClient,
	}
}

func (s *service) SayHello(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	req := &schema.SayHelloRequest{}

	err := json.Unmarshal([]byte(request.Body), req)
	if err != nil {
		return lambdaresponses.Respond400(fmt.Errorf("failed to unmarshal body"))
	}

	if req.Name == "Waldo" {
		return lambdaresponses.Respond400(fmt.Errorf("cannot use name Waldo!"))
	}

	message := fmt.Sprintf("Hello %s", req.Name)
	return lambdaresponses.Respond200(schema.SayHelloResponse{Message: message})
}

func (s *service) HostRoom(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if s.authClient == nil || s.ddbrepository == nil {
		log.Errorf("authClient or ddbrepository is nil")
		return lambdaresponses.Respond500()
	}

	req := &schema.HostRoomRequest{}

	err := json.Unmarshal([]byte(request.Body), req)
	if err != nil {
		return lambdaresponses.Respond400(fmt.Errorf("failed to unmarshal body"))
	}

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

func (s *service) FindRoom(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if s.ddbrepository == nil {
		log.Errorf("ddbrepository is nil")
		return lambdaresponses.Respond500()
	}

	req := &schema.FindRoomRequest{}

	err := json.Unmarshal([]byte(request.Body), req)
	if err != nil {
		return lambdaresponses.Respond400(fmt.Errorf("failed to unmarshal body"))
	}

	room, err := s.ddbrepository.FindRoom(ctx, req.ID)
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

func (s *service) JoinRoom(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if s.ddbrepository == nil || s.snsClient == nil || s.authClient == nil {

		log.Errorf("ddbrepository or snsClient or authClient is nil")
		return lambdaresponses.Respond500()
	}

	req := &schema.JoinRoomRequest{}

	err := json.Unmarshal([]byte(request.Body), req)
	if err != nil {
		return lambdaresponses.Respond400(fmt.Errorf("failed to unmarshal body"))
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

func (s *service) CastVote(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if s.snsClient == nil || s.authClient == nil {
		log.Errorf("snsClient or authClient is nil")
		return lambdaresponses.Respond500()
	}

	req := &schema.CastVoteRequest{}

	tokens := strings.Split(request.Headers["Authorization"], " ")
	tokenString := tokens[1]

	claims, err := s.authClient.GetClaims(tokenString)
	if err != nil {
		log.Errorf("failed to get claims")
	}

	err = json.Unmarshal([]byte(request.Body), req)
	if err != nil {
		return lambdaresponses.Respond400(fmt.Errorf("failed to unmarshal body"))
	}

	msg := schema.ParticipantVotedMessage{
		RoomID:          claims.RoomID,
		ParticipantName: claims.Name,
		Vote:            req.Vote,
	}

	err = s.snsClient.Publish(ctx, schema.ParticipantVoted, msg)
	if err != nil {
		log.Errorf("error publishing participant voted to sns: %w", err)
		return lambdaresponses.Respond500()
	}

	res := schema.CastVoteResponse{}

	return lambdaresponses.Respond200(res)
}
