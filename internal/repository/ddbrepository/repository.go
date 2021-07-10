package ddbrepository

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	awsDynamodb "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/jponc/estimatex-serverless/internal/types"
	"github.com/jponc/estimatex-serverless/pkg/dynamodb"
)

type Repository struct {
	dynamodbClient *dynamodb.Client
}

type roomItem struct {
	PK   string     `json:"PK"`
	SK   string     `json:"SK"`
	Data types.Room `json:"Data"`
}

type participantItem struct {
	PK   string            `json:"PK"`
	SK   string            `json:"SK"`
	Data types.Participant `json:"Data"`
}

// NewClient instantiates a repository
func NewClient(dynamodbClient *dynamodb.Client) (*Repository, error) {
	r := &Repository{
		dynamodbClient: dynamodbClient,
	}

	return r, nil
}

func (r *Repository) CreateRoom(ctx context.Context) (*types.Room, error) {
	roomID, err := r.generateRoomID(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to generate room ID (%w)", err)
	}

	room := &types.Room{
		ID:        roomID,
		CreatedAt: time.Now(),
	}

	item := struct {
		PK   string
		SK   string
		Data *types.Room
	}{
		PK:   fmt.Sprintf("Room_%s", room.ID),
		SK:   "RoomInfo",
		Data: room,
	}

	itemMap, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return nil, fmt.Errorf("failed to ddb marshal result item record, %v", err)
	}

	input := &awsDynamodb.PutItemInput{
		Item:      itemMap,
		TableName: aws.String(r.dynamodbClient.GetTableName()),
	}

	_, err = r.dynamodbClient.PutItem(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to put Room: %v", err)
	}

	return room, nil
}

func (r *Repository) CastVote(ctx context.Context, participant *types.Participant, vote string) error {
	participant.LatestVote = vote

	item := struct {
		PK   string
		SK   string
		Data *types.Participant
	}{
		PK:   fmt.Sprintf("Room_%s", participant.RoomID),
		SK:   fmt.Sprintf("Participant_%s", participant.Name),
		Data: participant,
	}

	itemMap, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("failed to ddb marshal result item record, %v", err)
	}

	input := &awsDynamodb.PutItemInput{
		Item:      itemMap,
		TableName: aws.String(r.dynamodbClient.GetTableName()),
	}

	_, err = r.dynamodbClient.PutItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to put Participant: %v", err)
	}

	return nil
}

func (r *Repository) CreateParticipant(ctx context.Context, roomID string, name string, isAdmin bool) (*types.Participant, error) {
	participant := &types.Participant{
		RoomID:    roomID,
		Name:      name,
		IsAdmin:   isAdmin,
		CreatedAt: time.Now(),
	}

	item := struct {
		PK   string
		SK   string
		Data *types.Participant
	}{
		PK:   fmt.Sprintf("Room_%s", participant.RoomID),
		SK:   fmt.Sprintf("Participant_%s", participant.Name),
		Data: participant,
	}

	itemMap, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return nil, fmt.Errorf("failed to ddb marshal result item record, %v", err)
	}

	input := &awsDynamodb.PutItemInput{
		Item:      itemMap,
		TableName: aws.String(r.dynamodbClient.GetTableName()),
	}

	_, err = r.dynamodbClient.PutItem(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to put Participant: %v", err)
	}

	return participant, nil
}

func (r *Repository) FindRoom(ctx context.Context, roomID string) (*types.Room, error) {
	i := roomItem{}

	input := &awsDynamodb.GetItemInput{
		Key: map[string]*awsDynamodb.AttributeValue{
			"PK": {
				S: aws.String(fmt.Sprintf("Room_%s", roomID)),
			},
			"SK": {
				S: aws.String("RoomInfo"),
			},
		},
		TableName: aws.String(r.dynamodbClient.GetTableName()),
	}

	output, err := r.dynamodbClient.GetItem(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to query room: %v", err)
	}

	if output.Item == nil {
		return nil, ErrNotFound
	}

	err = dynamodbattribute.UnmarshalMap(output.Item, &i)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal map: %v", err)
	}

	return &i.Data, nil
}

func (r *Repository) FindParticipant(ctx context.Context, roomID, participantName string) (*types.Participant, error) {
	i := participantItem{}

	input := &awsDynamodb.GetItemInput{
		Key: map[string]*awsDynamodb.AttributeValue{
			"PK": {
				S: aws.String(fmt.Sprintf("Room_%s", roomID)),
			},
			"SK": {
				S: aws.String(fmt.Sprintf("Participant_%s", participantName)),
			},
		},
		TableName: aws.String(r.dynamodbClient.GetTableName()),
	}

	output, err := r.dynamodbClient.GetItem(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to query participant: %v", err)
	}

	if output.Item == nil {
		return nil, ErrNotFound
	}

	err = dynamodbattribute.UnmarshalMap(output.Item, &i)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal map: %v", err)
	}

	return &i.Data, nil
}

func (r *Repository) FindParticipants(ctx context.Context, roomID string) (*[]types.Participant, error) {
	items := []participantItem{}

	input := &awsDynamodb.QueryInput{
		KeyConditionExpression: aws.String("PK = :PK and begins_with(SK, :SK)"),
		ExpressionAttributeValues: map[string]*awsDynamodb.AttributeValue{
			":PK": {
				S: aws.String(fmt.Sprintf("Room_%s", roomID)),
			},
			":SK": {
				S: aws.String("Participant_"),
			},
		},
		TableName: aws.String(r.dynamodbClient.GetTableName()),
	}

	output, err := r.dynamodbClient.Query(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to query participant: %v", err)
	}

	err = dynamodbattribute.UnmarshalListOfMaps(output.Items, &items)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal map: %v", err)
	}

	participants := []types.Participant{}
	for _, i := range items {
		participants = append(participants, i.Data)
	}

	return &participants, nil
}

func (r *Repository) generateRoomID(ctx context.Context) (string, error) {
	for {
		rand.Seed(time.Now().UnixNano())

		numCount := 6
		letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

		b := make([]rune, numCount)
		for i := range b {
			b[i] = letters[rand.Intn(len(letters))]
		}
		id := string(b)

		_, err := r.FindRoom(ctx, id)
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				return id, nil
			} else {
				return "", err
			}
		}
	}
}
