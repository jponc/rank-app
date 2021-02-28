package zenserp

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	pkgHttp "github.com/jponc/rank-app/pkg/http"
)

// Client interface
type Client interface {
	// Search provides zenserp search functionality
	Search(ctx context.Context, query string) (*QueryResult, error)
}

type client struct {
	apiKey     string
	baseURL    *url.URL
	httpClient *http.Client
}

// NewClient instantiates a zenserp client
func NewClient(apiKey string) (Client, error) {
	baseURL, err := url.Parse(zenserpBaseURL)
	if err != nil {
		return nil, fmt.Errorf("error parsing zenser Base URL (%w)", err)
	}

	return &client{
		apiKey:     apiKey,
		baseURL:    baseURL,
		httpClient: pkgHttp.DefaultHTTPClient(time.Duration(10 * time.Second)),
	}, nil
}

func (c *client) Search(ctx context.Context, query string) (*QueryResult, error) {
	res := &QueryResult{}
	endpoint := fmt.Sprintf(searchPath, query)
	err := c.getJSON(ctx, endpoint, res)

	if err != nil {
		return nil, fmt.Errorf("failed to query Zenserp (%s): %w", endpoint, err)
	}

	return res, nil
}
