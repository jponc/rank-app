package api

import "github.com/jponc/rank-app/internal/types"

const (
	ProcessKeyword     string = "ProcessKeyword"
	CrawlResultCreated string = "CrawlResultCreated"
)

type ProcessKeywordMessage struct {
	Keyword      string `json:"keyword"`
	SearchEngine string `json:"search_engine"`
	Device       string `json:"device"`
	Count        int    `json:"count"`
}

type CrawlResultCreatedMessage struct {
	CrawlResult types.CrawlResult `json:"crawl_result"`
}
