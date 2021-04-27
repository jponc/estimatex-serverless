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

type Repository interface {
	// CreateRoom creates Room
	CreateRoom(ctx context.Context) (*types.Room, error)
	// CreateParticipant creates Participant
	CreateParticipant(ctx context.Context, roomID string, name string, isAdmin bool) (*types.Participant, error)
	// FindRoom finds the room
	FindRoom(ctx context.Context, roomID string) (*types.Room, error)
	// FindParticipant finds the participant in the room
	FindParticipant(ctx context.Context, roomID, participantName string) (*types.Participant, error)
}

type repository struct {
	dynamodbClient dynamodb.Client
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
func NewClient(dynamodbClient dynamodb.Client) (Repository, error) {
	r := &repository{
		dynamodbClient: dynamodbClient,
	}

	return r, nil
}

func (r *repository) CreateRoom(ctx context.Context) (*types.Room, error) {
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

func (r *repository) CreateParticipant(ctx context.Context, roomID string, name string, isAdmin bool) (*types.Participant, error) {
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

func (r *repository) FindRoom(ctx context.Context, roomID string) (*types.Room, error) {
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

func (r *repository) FindParticipant(ctx context.Context, roomID, participantName string) (*types.Participant, error) {
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

func (r *repository) generateRoomID(ctx context.Context) (string, error) {
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
