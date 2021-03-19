package main

import (
	"fmt"
	"os"
)

// Config
type Config struct {
	ElasticsearchURL string
	AWSRegion        string
}

// NewConfig initialises a new config
func NewConfig() (*Config, error) {
	elasticsearchUrl, err := getEnv("ELASTICSEARCH_URL")
	if err != nil {
		return nil, err
	}

	awsRegion, err := getEnv("AWS_REGION")
	if err != nil {
		return nil, err
	}

	return &Config{
		ElasticsearchURL: elasticsearchUrl,
		AWSRegion:        awsRegion,
	}, nil
}

func getEnv(key string) (string, error) {
	v := os.Getenv(key)

	if v == "" {
		return "", fmt.Errorf("%s environment variable missing", key)
	}

	return v, nil
}
