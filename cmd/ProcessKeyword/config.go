package main

import (
	"fmt"
	"os"
)

// Config
type Config struct {
	ZenserpApiKey string
	AWSRegion     string
	DBTableName   string
	SNSPrefix     string
}

// NewConfig initialises a new config
func NewConfig() (*Config, error) {
	zenserpApiKey, err := getEnv("ZENSERP_API_KEY")
	if err != nil {
		return nil, err
	}

	awsRegion, err := getEnv("AWS_REGION")
	if err != nil {
		return nil, err
	}

	dbTableName, err := getEnv("DB_TABLE_NAME")
	if err != nil {
		return nil, err
	}

	snsPrefix, err := getEnv("SNS_PREFIX")
	if err != nil {
		return nil, err
	}

	return &Config{
		ZenserpApiKey: zenserpApiKey,
		AWSRegion:     awsRegion,
		DBTableName:   dbTableName,
		SNSPrefix:     snsPrefix,
	}, nil
}

func getEnv(key string) (string, error) {
	v := os.Getenv(key)

	if v == "" {
		return "", fmt.Errorf("%s environment variable missing", key)
	}

	return v, nil
}
