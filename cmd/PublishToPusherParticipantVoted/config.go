package main

import (
	"fmt"
	"os"
)

// Config
type Config struct {
	PusherAppID   string
	PusherKey     string
	PusherSecret  string
	PusherCluster string
}

// NewConfig initialises a new config
func NewConfig() (*Config, error) {
	appID, err := getEnv("PUSHER_APP_ID")
	if err != nil {
		return nil, err
	}

	key, err := getEnv("PUSHER_KEY")
	if err != nil {
		return nil, err
	}

	secret, err := getEnv("PUSHER_SECRET")
	if err != nil {
		return nil, err
	}

	cluster, err := getEnv("PUSHER_CLUSTER")
	if err != nil {
		return nil, err
	}

	return &Config{
		PusherAppID:   appID,
		PusherKey:     key,
		PusherSecret:  secret,
		PusherCluster: cluster,
	}, nil
}

func getEnv(key string) (string, error) {
	v := os.Getenv(key)

	if v == "" {
		return "", fmt.Errorf("%s environment variable missing", key)
	}

	return v, nil
}
