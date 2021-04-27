package api

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jponc/estimatex-serverless/api/schema"
	"github.com/jponc/estimatex-serverless/pkg/lambdaresponses"
)

// Service interface implements functions available for this service
type Service interface {
	SayHello(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
}

type service struct {
}

// NewService instantiates a new service
func NewService() Service {
	return &service{}
}

func (s *service) SayHello(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	req := &schema.SayHelloRequest{}

	err := json.Unmarshal([]byte(request.Body), req)
	if err != nil {
		return lambdaresponses.Respond400(fmt.Errorf("failed to unmarshall body"))
	}

	if req.Name == "Waldo" {
		return lambdaresponses.Respond400(fmt.Errorf("cannot use name Waldo!"))
	}

	message := fmt.Sprintf("Hello %s", req.Name)
	return lambdaresponses.Respond200(schema.SayHelloResponse{Message: message})
}
