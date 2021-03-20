package ddbrepository

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	awsDynamodb "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/jponc/rank-app/internal/types"
	"github.com/jponc/rank-app/pkg/dynamodb"
)

type Repository interface {
	// CreateResultItem creates ResultItem and persist to db
	CreateResultItem(ctx context.Context, resultItem types.ResultItem) error
}

type repository struct {
	dynamodbClient dynamodb.Client
}

// NewClient instantiates a repository
func NewClient(dynamodbClient dynamodb.Client) (Repository, error) {
	r := &repository{
		dynamodbClient: dynamodbClient,
	}

	return r, nil
}

func (r *repository) CreateResultItem(ctx context.Context, resultItem types.ResultItem) error {
	item := struct {
		PK   string
		SK   string
		Data *types.ResultItem
	}{
		PK:   fmt.Sprintf("ResultItem_%s", resultItem.ID.String()),
		SK:   "ResultItemInfo",
		Data: &resultItem,
	}

	itemMap, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("failed to ddb marshal result item record, %v", err)
	}

	input := &awsDynamodb.PutItemInput{
		Item:      itemMap,
		TableName: aws.String(r.dynamodbClient.GetTableName()),
	}

	_, err = r.dynamodbClient.PutItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to put ResultItem: %v", err)
	}

	return nil
}
