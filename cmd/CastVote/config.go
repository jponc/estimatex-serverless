package main

import (
	"fmt"
	"os"
)

// Config
type Config struct {
	AWSRegion string
	JWTSecret string
	SNSPrefix string
}

// NewConfig initialises a new config
func NewConfig() (*Config, error) {
	awsRegion, err := getEnv("AWS_REGION")
	if err != nil {
		return nil, err
	}

	jwtSecret, err := getEnv("JWT_SECRET")
	if err != nil {
		return nil, err
	}

	snsPrefix, err := getEnv("SNS_PREFIX")
	if err != nil {
		return nil, err
	}

	return &Config{
		AWSRegion: awsRegion,
		JWTSecret: jwtSecret,
		SNSPrefix: snsPrefix,
	}, nil
}

func getEnv(key string) (string, error) {
	v := os.Getenv(key)

	if v == "" {
		return "", fmt.Errorf("%s environment variable missing", key)
	}

	return v, nil
}