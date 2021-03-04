package converter

import (
	"github.com/jponc/rank-app/internal/types"
	"github.com/jponc/rank-app/pkg/zenserp"
)

func ZenserpQueryResultToCrawlResult(zenserpQueryResult *zenserp.QueryResult) *types.CrawlResult {
	queryInfo := zenserpQueryResult.Query
	items := []types.ResultItem{}

	for _, item := range zenserpQueryResult.ResulItems {
		items = append(items, types.ResultItem{
			Position:    item.Position,
			Title:       item.Title,
			URL:         item.URL,
			Description: item.Description,
		})
	}

	return &types.CrawlResult{
		Query:        queryInfo.Query,
		SearchEngine: queryInfo.SearchEngine,
		Device:       queryInfo.Device,
		URL:          queryInfo.URL,
		Items:        items,
	}
}
