package types

import (
	"time"

	"github.com/gofrs/uuid"
)

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
