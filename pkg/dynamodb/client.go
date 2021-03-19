package dynamodb

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// Client interface
type Client interface {
	// GetTableName gets the table name
	GetTableName() string
	// Scan scans the dynamodb table
	Scan(input *dynamodb.ScanInput) (*dynamodb.ScanOutput, error)
	// PutItem puts a new item in dynamodb table
	PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error)
	// BatchWriteItem writes multiple items to dynamodb table
	BatchWriteItem(input *dynamodb.BatchWriteItemInput) (*dynamodb.BatchWriteItemOutput, error)
	// Query queries dynamodb based on keys
	Query(input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error)
}

type client struct {
	dynamodbClient *dynamodb.DynamoDB
	tableName      string
}

// NewClient instantiates a DynamoDB Client
func NewClient(awsRegion, tableName string) (Client, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
	})

	if err != nil {
		return nil, fmt.Errorf("cannot create aws session: %v", err)
	}

	dynamodbClient := dynamodb.New(sess)

	c := &client{
		dynamodbClient: dynamodbClient,
		tableName:      tableName,
	}

	return c, nil
}

func (c *client) Scan(input *dynamodb.ScanInput) (*dynamodb.ScanOutput, error) {
	return c.dynamodbClient.Scan(input)
}

func (c *client) PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	return c.dynamodbClient.PutItem(input)
}

func (c *client) BatchWriteItem(input *dynamodb.BatchWriteItemInput) (*dynamodb.BatchWriteItemOutput, error) {
	return c.dynamodbClient.BatchWriteItem(input)
}

func (c *client) Query(input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
	return c.dynamodbClient.Query(input)
}

func (c *client) GetTableName() string {
	return c.tableName
}
