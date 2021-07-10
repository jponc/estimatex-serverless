package pusher

import (
	"context"

	push "github.com/pusher/pusher-http-go"
)

type Client struct {
	pusherClient *push.Client
}

// NewClient instantiates a DynamoDB Client
func NewClient(appID, key, secret, cluster string) (*Client, error) {
	pusherClient := push.Client{
		AppID:   appID,
		Key:     key,
		Secret:  secret,
		Cluster: cluster,
		Secure:  true,
	}

	c := &Client{
		pusherClient: &pusherClient,
	}

	return c, nil
}

func (c *Client) Trigger(ctx context.Context, channel, eventName string, data interface{}) error {
	return c.pusherClient.Trigger(channel, eventName, data)
}
