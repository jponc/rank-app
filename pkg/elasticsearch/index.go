package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
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

func (c *client) DeleteIndex(ctx context.Context, index string) error {
	req := esapi.IndicesDeleteRequest{
		Index: []string{index},
	}

	res, err := req.Do(ctx, c.esClient)
	if err != nil || res.IsError() {
		return fmt.Errorf("error deleting index %s: %v", index, err)
	}

	return nil
}

func (c *client) Index(ctx context.Context, index, id string, body interface{}) error {
	bodyByte, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("index cannot marshal body: %v, err: %v", body, err)
	}

	req := esapi.IndexRequest{
		Index:      index,
		DocumentID: id,
		Body:       bytes.NewReader(bodyByte),
	}

	if res, err := req.Do(ctx, c.esClient); err != nil || res.IsError() {
		return fmt.Errorf("error es index: %v", err)
	}

	return nil
}
