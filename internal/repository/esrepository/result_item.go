package esrepository

import (
	"context"
	"fmt"

	"github.com/jponc/rank-app/internal/types"
)

func (r *repository) IndexResultItem(ctx context.Context, resultItem types.ResultItem) error {
	esResultItem := &types.ESResultItem{}

	if err := esResultItem.Unmarshal(resultItem); err != nil {
		return fmt.Errorf("failed indexing result item: %v", err)
	}

	return nil

}
