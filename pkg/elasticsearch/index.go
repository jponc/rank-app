package elasticsearch

import (
	"context"
	"fmt"

	"github.com/elastic/go-elasticsearch/v7/esapi"
)

func (c *client) IndexExists(ctx context.Context, index string) (*bool, error) {
	req := esapi.IndicesExistsRequest{
		Index: []string{index},
	}

	res, err := req.Do(ctx, c.esClient)
	if err != nil {
		return nil, fmt.Errorf("error getting index exists: %v", err)
	}

	exists := !res.IsError()
	return &exists, nil
}

func (c *client) CreateIndex(ctx context.Context, index string) error {
	req := esapi.IndicesCreateRequest{
		Index: index,
	}

	res, err := req.Do(ctx, c.esClient)
	if err != nil || res.IsError() {
		return fmt.Errorf("error creating index %s: %v", index, err)
	}

	return nil
}
