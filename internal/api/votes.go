package api

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jponc/estimatex-serverless/api/schema"
	"github.com/jponc/estimatex-serverless/pkg/lambdaresponses"
	log "github.com/sirupsen/logrus"
)

func (s *Service) CastVote(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if s.snsClient == nil {
		log.Errorf("snsClient or authClient is nil")
		return lambdaresponses.Respond500()
	}

	req := &schema.CastVoteRequest{}

	roomID, ok := request.RequestContext.Authorizer["RoomID"].(string)
	if !ok {
		return lambdaresponses.Respond500()
	}

	name, ok := request.RequestContext.Authorizer["Name"].(string)
	if !ok {
		return lambdaresponses.Respond500()
	}

	err := json.Unmarshal([]byte(request.Body), req)
	if err != nil {
		return lambdaresponses.Respond400(fmt.Errorf("failed to unmarshal body"))
	}

	if req.Vote == "" {
		return lambdaresponses.Respond400(fmt.Errorf("vote can't be blank"))
	}

	log.Infof("roomID: %s, name: %s, vote: %s", roomID, name, req.Vote)
	err = s.ddbrepository.CastVote(ctx, roomID, name, req.Vote)
	if err != nil {
		return lambdaresponses.Respond500()
	}

	msg := schema.ParticipantVotedMessage{
		RoomID:          roomID,
		ParticipantName: name,
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

func (s *Service) RevealVotes(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if s.snsClient == nil {
		log.Errorf("snsClient or authClient is nil")
		return lambdaresponses.Respond500()
	}

	isAdmin, ok := request.RequestContext.Authorizer["IsAdmin"].(string)
	if !ok {
		log.Errorf("no is admin")
		return lambdaresponses.Respond500()
	}

	if isAdmin != "true" {
		return lambdaresponses.Respond400(fmt.Errorf("not allowed"))
	}

	roomID, ok := request.RequestContext.Authorizer["RoomID"].(string)
	if !ok {
		log.Errorf("no room id")
		return lambdaresponses.Respond500()
	}

	msg := schema.RevealVotesMessage{
		RoomID: roomID,
	}

	err := s.snsClient.Publish(ctx, schema.RevealVotes, msg)
	if err != nil {
		log.Errorf("error doing RevealVotes to sns: %w", err)
		return lambdaresponses.Respond500()
	}

	res := schema.RevealVotesResponse{}

	return lambdaresponses.Respond200(res)
}

func (s *Service) ResetVotes(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if s.snsClient == nil {
		log.Errorf("snsClient or authClient is nil")
		return lambdaresponses.Respond500()
	}

	isAdmin, ok := request.RequestContext.Authorizer["IsAdmin"].(string)
	if !ok {
		log.Errorf("no is admin")
		return lambdaresponses.Respond500()
	}

	if isAdmin != "true" {
		return lambdaresponses.Respond400(fmt.Errorf("not allowed"))
	}

	roomID, ok := request.RequestContext.Authorizer["RoomID"].(string)
	if !ok {
		log.Errorf("no room id")
		return lambdaresponses.Respond500()
	}

	msg := schema.ResetVotesMessage{
		RoomID: roomID,
	}

	err := s.snsClient.Publish(ctx, schema.ResetVotes, msg)
	if err != nil {
		log.Errorf("error doing RevealVotes to sns: %w", err)
		return lambdaresponses.Respond500()
	}

	res := schema.ResetVotesResponse{}

	return lambdaresponses.Respond200(res)
}
