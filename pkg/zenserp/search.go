package zenserp

import (
	"context"
	"fmt"
)

func (c *client) Search(ctx context.Context, query string, num int) (*QueryResult, error) {
	if num > 100 {
		return nil, fmt.Errorf("result count (num) of %d, not allowed", num)
	}

	res := &QueryResult{}
	endpoint := fmt.Sprintf(searchPath, query, num)
	err := c.getJSON(ctx, endpoint, res)

	if err != nil {
		return nil, fmt.Errorf("failed to query Zenserp (%s): %w", endpoint, err)
	}

	return res, nil
}
