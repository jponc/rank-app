package repository

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	awsDynamodb "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gofrs/uuid"
	"github.com/jponc/rank-app/internal/converter"
	"github.com/jponc/rank-app/internal/types"
	"github.com/jponc/rank-app/pkg/dynamodb"
	"github.com/jponc/rank-app/pkg/zenserp"
)

type Repository interface {
	CreateCrawlResult(zenserpResult *zenserp.QueryResult) (*types.CrawlResult, error)
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

func (r *repository) CreateCrawlResult(zenserpQueryResult *zenserp.QueryResult) (*types.CrawlResult, error) {
	crawlResult := converter.ZenserpQueryResultToCrawlResult(zenserpQueryResult)
	crawlResult.ID = uuid.Must(uuid.NewV4())
	crawlResult.CreatedAt = time.Now()

	crawlResultMap, err := dynamodbattribute.MarshalMap(crawlResult)
	if err != nil {
		return nil, fmt.Errorf("failed to DynamoDB marshal Record, %v", err)
	}

	input := &awsDynamodb.PutItemInput{
		Item: map[string]*awsDynamodb.AttributeValue{
			"PK": {
				S: aws.String(fmt.Sprintf("CrawResult_%s", crawlResult.ID.String())),
			},
			"SK": {
				S: aws.String("CrawResultData"),
			},
			"Data": {
				M: crawlResultMap,
			},
		},
		TableName: aws.String(r.dynamodbClient.GetTableName()),
	}

	_, err = r.dynamodbClient.PutItem(input)
	if err != nil {
		return nil, fmt.Errorf("failed to put CrawlResult: %v", err)
	}

	return crawlResult, nil
}
