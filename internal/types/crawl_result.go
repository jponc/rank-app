package types

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jponc/rank-app/pkg/zenserp"
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

func (t *CrawlResult) Unmarshal(src interface{}) error {
	switch src.(type) {
	case *zenserp.QueryResult:
		return t.unmarshalZenserpQueryResult(*src.(*zenserp.QueryResult))
	}

	return fmt.Errorf("Failed to unmarshal types.CrawlResult: '%v'", src)
}

func (t *CrawlResult) unmarshalZenserpQueryResult(zenserpQueryResult zenserp.QueryResult) error {
	queryInfo := zenserpQueryResult.Query
	items := []ResultItem{}

	for _, item := range zenserpQueryResult.ResulItems {
		items = append(items, ResultItem{
			Position:    item.Position,
			Title:       item.Title,
			URL:         item.URL,
			Description: item.Description,
		})
	}

	t = &CrawlResult{
		Query:        queryInfo.Query,
		SearchEngine: queryInfo.SearchEngine,
		Device:       queryInfo.Device,
		URL:          queryInfo.URL,
		Items:        items,
	}

	return nil
}
