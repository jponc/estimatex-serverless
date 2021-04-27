package api

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jponc/estimatex-serverless/api/schema"
	"github.com/jponc/estimatex-serverless/internal/auth"
	"github.com/jponc/estimatex-serverless/internal/repository/ddbrepository"
	"github.com/jponc/estimatex-serverless/pkg/lambdaresponses"
	"github.com/jponc/estimatex-serverless/pkg/sns"
)

// Service interface implements functions available for this service
type Service interface {
	// SayHello is a placeholder endpoint
	SayHello(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
	// HostRoom creates a new room
	HostRoom(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
	// FindRoom finds the room given a room ID
	FindRoom(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
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
		return lambdaresponses.Respond500()
	}

	req := &schema.HostRoomRequest{}

	err := json.Unmarshal([]byte(request.Body), req)
	if err != nil {
		return lambdaresponses.Respond400(fmt.Errorf("failed to unmarshal body"))
	}

	room, err := s.ddbrepository.CreateRoom(ctx)
	if err != nil {
		return lambdaresponses.Respond500()
	}

	participant, err := s.ddbrepository.CreateParticipant(ctx, room.ID, req.Name, true)
	if err != nil {
		return lambdaresponses.Respond500()
	}

	token, err := s.authClient.CreateAccessToken(*participant)
	if err != nil {
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
		return lambdaresponses.Respond500()
	}

	req := &schema.FindRoomRequest{}

	err := json.Unmarshal([]byte(request.Body), req)
	if err != nil {
		return lambdaresponses.Respond400(fmt.Errorf("failed to unmarshal body"))
	}

	room, err := s.ddbrepository.FindRoom(ctx, req.ID)
	if err != nil {
		return lambdaresponses.Respond500()
	}
	if room == nil {
		return lambdaresponses.Respond404(fmt.Errorf("room not found"))
	}

	res := schema.FindRoomResponse{
		Room: *room,
	}

	return lambdaresponses.Respond200(res)
}
