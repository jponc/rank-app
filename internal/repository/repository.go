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
	log "github.com/sirupsen/logrus"
)

type Repository interface {
	// CreateCrawlResult creates CrawlResult and persist to DB
	CreateCrawlResult(zenserpResult *zenserp.QueryResult) (*types.CrawlResult, error)
	// AddCrawlResultToLatest adds crawl result to LatestTemporary
	AddCrawlResultToLatest(result *types.CrawlResult) error
	// GetLatestCrawlResults gets all the latest crawl results
	GetLatestCrawlResults() (*[]types.CrawlResult, error)
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

	item := struct {
		PK   string
		SK   string
		Data *types.CrawlResult
	}{
		PK:   fmt.Sprintf("CrawlResult_%s", crawlResult.ID.String()),
		SK:   "CrawlResult",
		Data: crawlResult,
	}

	itemMap, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return nil, fmt.Errorf("failed to DynamoDB marshal Record, %v", err)
	}

	input := &awsDynamodb.PutItemInput{
		Item:      itemMap,
		TableName: aws.String(r.dynamodbClient.GetTableName()),
	}

	_, err = r.dynamodbClient.PutItem(input)
	if err != nil {
		return nil, fmt.Errorf("failed to put CrawlResult: %v", err)
	}

	return crawlResult, nil
}

func (r *repository) AddCrawlResultToLatest(crawlResult *types.CrawlResult) error {
	item := struct {
		PK   string
		SK   string
		Data *types.CrawlResult
	}{
		PK:   "LatestCrawlResults",
		SK:   fmt.Sprintf("CrawlResult_%s_%s_%s", crawlResult.Query, crawlResult.SearchEngine, crawlResult.Device),
		Data: crawlResult,
	}

	itemMap, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("failed to DynamoDB marshal Record, %v", err)
	}

	input := &awsDynamodb.PutItemInput{
		Item:      itemMap,
		TableName: aws.String(r.dynamodbClient.GetTableName()),
	}

	_, err = r.dynamodbClient.PutItem(input)
	if err != nil {
		return fmt.Errorf("failed to put LatestCrawlResults: %v", err)
	}

	return nil
}

func (r *repository) GetLatestCrawlResults() (*[]types.CrawlResult, error) {
	values := map[string]string{
		":pk": "LatestCrawlResults",
		":sk": "CrawlResult_",
	}
	valuesMap, err := dynamodbattribute.MarshalMap(values)
	if err != nil {
		return nil, fmt.Errorf("failed to DynamoDB marshal Record, %v", err)
	}

	input := &awsDynamodb.QueryInput{
		TableName:              aws.String(r.dynamodbClient.GetTableName()),
		KeyConditionExpression: aws.String("#pk = :pk and begins_with(#sk, :sk)"),
		ExpressionAttributeNames: map[string]*string{
			"#pk": aws.String("PK"),
			"#sk": aws.String("SK"),
		},
		ExpressionAttributeValues: valuesMap,
	}

	res, err := r.dynamodbClient.Query(input)
	if err != nil {
		return nil, fmt.Errorf("failed to query latest crawl results: %v", err)
	}

	type item struct {
		PK   string
		SK   string
		Data types.CrawlResult
	}

	crawlResults := []types.CrawlResult{}
	log.Infof("Items length: %d", len(res.Items))

	for _, resItem := range res.Items {
		i := &item{}

		err = dynamodbattribute.UnmarshalMap(resItem, &i)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal result item %v", err)
		}

		crawlResults = append(crawlResults, i.Data)
	}

	return &crawlResults, nil
}
