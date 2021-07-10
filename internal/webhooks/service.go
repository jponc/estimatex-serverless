package webhooks

import (
	"context"
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jponc/estimatex-serverless/api/schema"
	"github.com/jponc/estimatex-serverless/pkg/pusher"
)

type Service struct {
	pusherClient *pusher.Client
}

// NewService instantiates a new service
func NewService(pusherClient *pusher.Client) *Service {
	return &Service{
		pusherClient: pusherClient,
	}
}

func (s *Service) PublishToPusherParticipantJoined(ctx context.Context, snsEvent events.SNSEvent) {
	snsMsg := snsEvent.Records[0].SNS.Message

	var msg schema.ParticipantJoinedMessage
	err := json.Unmarshal([]byte(snsMsg), &msg)
	if err != nil {
		log.Fatalf("unable to unarmarshal message: %v", err)
	}

	if s.pusherClient == nil {
		log.Fatalf("pusherClient not defined")
	}

	channel := fmt.Sprintf("room-%s", msg.RoomID)
	event := "participant-joined"
	data := map[string]string{
		"room_id":          msg.RoomID,
		"participant_name": msg.ParticipantName,
	}

	err = s.pusherClient.Trigger(ctx, channel, event, data)
	if err != nil {
		log.Fatalf("failed to trigger push: %w", err)
	}
}

func (s *Service) PublishToPusherParticipantVoted(ctx context.Context, snsEvent events.SNSEvent) {
	snsMsg := snsEvent.Records[0].SNS.Message

	var msg schema.ParticipantVotedMessage
	err := json.Unmarshal([]byte(snsMsg), &msg)
	if err != nil {
		log.Fatalf("unable to unarmarshal message: %v", err)
	}

	if s.pusherClient == nil {
		log.Fatalf("pusherClient not defined")
	}

	channel := fmt.Sprintf("room-%s", msg.RoomID)
	event := "participant-voted"
	data := map[string]string{
		"room_id":          msg.RoomID,
		"participant_name": msg.ParticipantName,
		"vote":             msg.Vote,
	}

	err = s.pusherClient.Trigger(ctx, channel, event, data)
	if err != nil {
		log.Fatalf("failed to trigger push: %w", err)
	}
}

func (s *Service) PublishToPusherResetVotes(ctx context.Context, snsEvent events.SNSEvent) {
	snsMsg := snsEvent.Records[0].SNS.Message

	var msg schema.ParticipantVotedMessage
	err := json.Unmarshal([]byte(snsMsg), &msg)
	if err != nil {
		log.Fatalf("unable to unarmarshal message: %v", err)
	}

	if s.pusherClient == nil {
		log.Fatalf("pusherClient not defined")
	}

	channel := fmt.Sprintf("room-%s", msg.RoomID)
	event := "reset-votes"
	data := map[string]string{
		"room_id": msg.RoomID,
	}

	err = s.pusherClient.Trigger(ctx, channel, event, data)
	if err != nil {
		log.Fatalf("failed to trigger push: %w", err)
	}
}

func (s *Service) PublishToPusherRevealVotes(ctx context.Context, snsEvent events.SNSEvent) {
	snsMsg := snsEvent.Records[0].SNS.Message

	var msg schema.ParticipantVotedMessage
	err := json.Unmarshal([]byte(snsMsg), &msg)
	if err != nil {
		log.Fatalf("unable to unarmarshal message: %v", err)
	}

	if s.pusherClient == nil {
		log.Fatalf("pusherClient not defined")
	}

	channel := fmt.Sprintf("room-%s", msg.RoomID)
	event := "reveal-votes"
	data := map[string]string{
		"room_id": msg.RoomID,
	}

	err = s.pusherClient.Trigger(ctx, channel, event, data)
	if err != nil {
		log.Fatalf("failed to trigger push: %w", err)
	}
}
