package esrepository

import (
	"context"
	"fmt"
)

type AttributeConfig struct {
	Type string `json:"type"`
}

type ResultItemMappingProperties struct {
	ID                      AttributeConfig `json:"id"`
	CrawlResultID           AttributeConfig `json:"crawl_result_id"`
	CrawlResultQuery        AttributeConfig `json:"crawl_result_query"`
	CrawlResultSearchEngine AttributeConfig `json:"crawl_result_search_engine"`
	CrawlResultDevice       AttributeConfig `json:"crawl_result_device"`
	CrawlResultURL          AttributeConfig `json:"craw_result_url"`
	CrawlResultCreatedAt    AttributeConfig `json:"crawl_result_created_at"`
	ItemPosition            AttributeConfig `json:"item_position"`
	ItemTitle               AttributeConfig `json:"item_title"`
	ItemURL                 AttributeConfig `json:"item_url"`
	ItemDescription         AttributeConfig `json:"item_description"`
}

type ResultItemMapping struct {
	Properties ResultItemMappingProperties `json:"properties"`
}

func (r *repository) UpdateResultItemIndexMapping(ctx context.Context) error {
	body := ResultItemMapping{
		Properties: ResultItemMappingProperties{
			ID:                      AttributeConfig{Type: "keyword"},
			CrawlResultID:           AttributeConfig{Type: "keyword"},
			CrawlResultQuery:        AttributeConfig{Type: "keyword"},
			CrawlResultSearchEngine: AttributeConfig{Type: "keyword"},
			CrawlResultDevice:       AttributeConfig{Type: "keyword"},
			CrawlResultURL:          AttributeConfig{Type: "keyword"},
			CrawlResultCreatedAt:    AttributeConfig{Type: "date"},
			ItemPosition:            AttributeConfig{Type: "integer"},
			ItemTitle:               AttributeConfig{Type: "keyword"},
			ItemURL:                 AttributeConfig{Type: "keyword"},
			ItemDescription:         AttributeConfig{Type: "keyword"},
		},
	}

	err := r.esClient.UpdateMapping(ctx, resultItemIndexName, body)
	if err != nil {
		return fmt.Errorf("failed to update es mapping: %v", err)
	}

	return nil
}
