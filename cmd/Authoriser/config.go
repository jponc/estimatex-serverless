package main

import (
	"fmt"
	"os"
)

// Config
type Config struct {
	JWTSecret string
}

// NewConfig initialises a new config
func NewConfig() (*Config, error) {
	jwtSecret, err := getEnv("JWT_SECRET")
	if err != nil {
		return nil, err
	}

	return &Config{
		JWTSecret: jwtSecret,
	}, nil
}

func getEnv(key string) (string, error) {
	v := os.Getenv(key)

	if v == "" {
		return "", fmt.Errorf("%s environment variable missing", key)
	}

	return v, nil
}
