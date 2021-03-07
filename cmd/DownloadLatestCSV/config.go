package main

import (
	"fmt"
	"os"
)

// Config
type Config struct {
	AWSRegion    string
	S3BucketName string
}

// NewConfig initialises a new config
func NewConfig() (*Config, error) {
	awsRegion, err := getEnv("AWS_REGION")
	if err != nil {
		return nil, err
	}

	s3BucketName, err := getEnv("S3_BUCKET_NAME")
	if err != nil {
		return nil, err
	}

	return &Config{
		AWSRegion:    awsRegion,
		S3BucketName: s3BucketName,
	}, nil
}

func getEnv(key string) (string, error) {
	v := os.Getenv(key)

	if v == "" {
		return "", fmt.Errorf("%s environment variable missing", key)
	}

	return v, nil
}
