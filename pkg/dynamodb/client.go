package dynamodb

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-xray-sdk-go/xray"
)

// Client interface
type Client interface {
	// GetTableName gets the table name
	GetTableName() string
	// Scan scans the dynamodb table
	Scan(ctx context.Context, input *dynamodb.ScanInput) (*dynamodb.ScanOutput, error)
	// PutItem puts a new item in dynamodb table
	PutItem(ctx context.Context, input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error)
	// BatchWriteItem writes multiple items to dynamodb table
	BatchWriteItem(ctx context.Context, input *dynamodb.BatchWriteItemInput) (*dynamodb.BatchWriteItemOutput, error)
	// Query queries dynamodb based on keys
	Query(ctx context.Context, input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error)
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
	xray.AWS(dynamodbClient.Client)

	c := &client{
		dynamodbClient: dynamodbClient,
		tableName:      tableName,
	}

	return c, nil
}

func (c *client) Scan(ctx context.Context, input *dynamodb.ScanInput) (*dynamodb.ScanOutput, error) {
	return c.dynamodbClient.ScanWithContext(ctx, input)
}

func (c *client) PutItem(ctx context.Context, input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	return c.dynamodbClient.PutItemWithContext(ctx, input)
}

func (c *client) BatchWriteItem(ctx context.Context, input *dynamodb.BatchWriteItemInput) (*dynamodb.BatchWriteItemOutput, error) {
	return c.dynamodbClient.BatchWriteItemWithContext(ctx, input)
}

func (c *client) Query(ctx context.Context, input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
	return c.dynamodbClient.QueryWithContext(ctx, input)
}

func (c *client) GetTableName() string {
	return c.tableName
}
