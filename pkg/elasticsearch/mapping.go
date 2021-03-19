package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch/v7/esapi"
)

func (c *client) UpdateMapping(ctx context.Context, index string, body interface{}) error {

	bodyByte, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("indices put mapping cannot marshal body: %v, err: %v", body, err)
	}

	existsPtr, err := c.IndexExists(ctx, index)
	if err != nil {
		return err
	}

	if !*existsPtr {
		if err = c.CreateIndex(ctx, index); err != nil {
			return fmt.Errorf("can't update mapping: %v", err)
		}
	}

	req := esapi.IndicesPutMappingRequest{
		Index: []string{index},
		Body:  bytes.NewReader(bodyByte),
	}

	res, err := req.Do(ctx, c.esClient)
	if err != nil {
		return fmt.Errorf("error updating ES index mapping: %v", err)
	}

	fmt.Printf("successfully updated ES mapping: %v\n", res)

	return nil
}
