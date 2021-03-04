package types

import (
	"time"

	"github.com/gofrs/uuid"
)

type CrawlResult struct {
	ID           uuid.UUID    `json:"id"`
	Query        string       `json:"query"`
	SearchEngine string       `json:"search_engine"`
	Device       string       `json:"device"`
	URL          string       `json:"url"`
	Items        []ResultItem `json:"items"`
	CreatedAt    time.Time    `json:"created_at"`
}

type ResultItem struct {
	Position    int    `json:"position"`
	Title       string `json:"title"`
	URL         string `json:"url"`
	Description string `json:"description"`
}
