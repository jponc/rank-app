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

type ESResultItem struct {
	CrawlResultID           uuid.UUID `json:"crawl_result_id"`
	CrawlResultQuery        string    `json:"crawl_result_query"`
	CrawlResultSearchEngine string    `json:"crawl_result_search_engine"`
	CrawlResultDevice       string    `json:"crawl_result_device"`
	CrawlResultURL          string    `json:"craw_result_url"`
	CrawlResultCreatedAt    time.Time `json:"crawl_result_created_at"`
	ItemPosition            int       `json:"item_position"`
	ItemTitle               string    `json:"item_title"`
	ItemURL                 string    `json:"item_url"`
	ItemDescription         string    `json:"item_description"`
}
