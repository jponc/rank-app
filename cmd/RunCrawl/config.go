package main

import (
	"fmt"
	"os"
)

// Config holds all configuration related to zenserp
type Config struct {
	AWSRegion string
	SNSPrefix string
}

// NewConfig initialises a new config
func NewConfig() (*Config, error) {
	snsPrefix, err := getEnv("SNS_PREFIX")
	if err != nil {
		return nil, err
	}

	awsRegion, err := getEnv("AWS_REGION")
	if err != nil {
		return nil, err
	}

	return &Config{
		AWSRegion: awsRegion,
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
