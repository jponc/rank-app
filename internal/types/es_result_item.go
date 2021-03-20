package types

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"
)

type ESResultItem struct {
	ID            uuid.UUID `json:"id"`
	CrawlResultID uuid.UUID `json:"crawl_result_id"`
	Query         string    `json:"query"`
	SearchEngine  string    `json:"search_engine"`
	Device        string    `json:"device"`
	QueryURL      string    `json:"query_url"`
	Position      int       `json:"position"`
	Title         string    `json:"title"`
	ItemURL       string    `json:"item_url"`
	Description   string    `json:"description"`
	CreatedAt     time.Time `json:"created_at"`
}

func (t *ESResultItem) Unmarshal(src interface{}) error {
	switch src.(type) {
	case *ResultItem:
		return t.unmarshalResultItem(*src.(*ResultItem))
	case ResultItem:
		return t.unmarshalResultItem(src.(ResultItem))
	}

	return fmt.Errorf("Failed to unmarshal types.ESResultItem: '%v'", src)
}

func (t *ESResultItem) unmarshalResultItem(resultItem ResultItem) error {
	t = &ESResultItem{
		ID:            resultItem.ID,
		CrawlResultID: resultItem.CrawlResultID,
		Query:         resultItem.Query,
		SearchEngine:  resultItem.SearchEngine,
		Device:        resultItem.Device,
		QueryURL:      resultItem.QueryURL,
		Position:      resultItem.Position,
		Title:         resultItem.Title,
		ItemURL:       resultItem.ItemURL,
		Description:   resultItem.Description,
		CreatedAt:     resultItem.CreatedAt,
	}

	return nil
}
