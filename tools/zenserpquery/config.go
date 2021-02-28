package main

import (
	"fmt"
	"os"
)

// Config holds all configuration related to zenserp
type Config struct {
	APIKey string
}

// NewConfig initialises a new config
func NewConfig() (*Config, error) {
	apiKey, err := getEnv("API_KEY")
	if err != nil {
		return nil, err
	}

	return &Config{
		APIKey: apiKey,
	}, nil
}

func getEnv(key string) (string, error) {
	v := os.Getenv(key)

	if v == "" {
		return "", fmt.Errorf("%s environment variable missing", key)
	}

	return v, nil
}
