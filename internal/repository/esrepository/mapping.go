package esrepository

import (
	"context"
	"fmt"
)

type AttributeConfig struct {
	Type string `json:"type"`
}

type ResultItemMappingProperties struct {
	ID            AttributeConfig `json:"id"`
	CrawlResultID AttributeConfig `json:"crawl_result_id"`
	Query         AttributeConfig `json:"query"`
	SearchEngine  AttributeConfig `json:"search_engine"`
	Device        AttributeConfig `json:"device"`
	QueryURL      AttributeConfig `json:"query_url"`
	Position      AttributeConfig `json:"position"`
	Title         AttributeConfig `json:"title"`
	ItemURL       AttributeConfig `json:"item_url"`
	Description   AttributeConfig `json:"description"`
	CreatedAt     AttributeConfig `json:"created_at"`
}

type ResultItemMapping struct {
	Properties ResultItemMappingProperties `json:"properties"`
}

func (r *repository) UpdateResultItemIndexMapping(ctx context.Context) error {
	body := ResultItemMapping{
		Properties: ResultItemMappingProperties{
			ID:            AttributeConfig{Type: "keyword"},
			CrawlResultID: AttributeConfig{Type: "keyword"},
			Query:         AttributeConfig{Type: "keyword"},
			SearchEngine:  AttributeConfig{Type: "keyword"},
			Device:        AttributeConfig{Type: "keyword"},
			QueryURL:      AttributeConfig{Type: "keyword"},
			Position:      AttributeConfig{Type: "integer"},
			Title:         AttributeConfig{Type: "keyword"},
			ItemURL:       AttributeConfig{Type: "keyword"},
			Description:   AttributeConfig{Type: "keyword"},
			CreatedAt:     AttributeConfig{Type: "date"},
		},
	}

	err := r.esClient.UpdateMapping(ctx, resultItemIndexName, body)
	if err != nil {
		return fmt.Errorf("failed to update es mapping: %v", err)
	}

	return nil
}
