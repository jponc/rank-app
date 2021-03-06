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
	// CreateCrawlResult creates CrawlResult and persist to DB
	CreateCrawlResult(zenserpResult *zenserp.QueryResult) (*types.CrawlResult, error)
	// AddCrawlResultToLatest adds crawl result to LatestTemporary
	AddCrawlResultToLatest(result *types.CrawlResult) error
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

func (r *repository) AddCrawlResultToLatest(crawlResult *types.CrawlResult) error {
	crawlResultMap, err := dynamodbattribute.MarshalMap(crawlResult)
	if err != nil {
		return fmt.Errorf("failed to DynamoDB marshal Record, %v", err)
	}

	sk := fmt.Sprintf("CrawlResult_%s_%s_%s", crawlResult.Query, crawlResult.SearchEngine, crawlResult.Device)

	input := &awsDynamodb.PutItemInput{
		Item: map[string]*awsDynamodb.AttributeValue{
			"PK": {
				S: aws.String("LatestCrawlResults"),
			},
			"SK": {
				S: aws.String(sk),
			},
			"Data": {
				M: crawlResultMap,
			},
		},
		TableName: aws.String(r.dynamodbClient.GetTableName()),
	}

	_, err = r.dynamodbClient.PutItem(input)
	if err != nil {
		return fmt.Errorf("failed to put LatestCrawlResults: %v", err)
	}

	return nil
}
