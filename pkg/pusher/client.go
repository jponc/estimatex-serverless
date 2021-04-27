package pusher

import (
	"context"

	push "github.com/pusher/pusher-http-go"
)

// Client interface
type Client interface {
	// Trigger triggers a new event to a pusher channel with corresponding data
	Trigger(ctx context.Context, channel, eventName string, data interface{}) error
}

type client struct {
	pusherClient *push.Client
}

// NewClient instantiates a DynamoDB Client
func NewClient(appID, key, secret, cluster string) (Client, error) {
	pusherClient := push.Client{
		AppID:   appID,
		Key:     key,
		Secret:  secret,
		Cluster: cluster,
		Secure:  true,
	}

	c := &client{
		pusherClient: &pusherClient,
	}

	return c, nil
}

func (c *client) Trigger(ctx context.Context, channel, eventName string, data interface{}) error {
	return c.pusherClient.Trigger(channel, eventName, data)
}
