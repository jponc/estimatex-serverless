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

type Service struct {
	ddbrepository *ddbrepository.Repository
	snsClient     *sns.Client
	authClient    *auth.Client
}

// NewService instantiates a new service
func NewService(ddbrepository *ddbrepository.Repository, snsClient *sns.Client, authClient *auth.Client) *Service {
	return &Service{
		ddbrepository: ddbrepository,
		snsClient:     snsClient,
		authClient:    authClient,
	}
}

func (s *Service) SayHello(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
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
