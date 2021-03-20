package api

import "github.com/jponc/rank-app/internal/types"

const (
	ProcessKeyword    string = "ProcessKeyword"
	ResultItemCreated string = "ResultItemCreated"
)

type ProcessKeywordMessage struct {
	Keyword      string `json:"keyword"`
	SearchEngine string `json:"search_engine"`
	Device       string `json:"device"`
	Count        int    `json:"count"`
}

type ResultItemCreatedMessage struct {
	ResultItem types.ResultItem `json:"result_item"`
}
