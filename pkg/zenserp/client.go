package zenserp

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	pkgHttp "github.com/jponc/rank-app/pkg/http"
)

// Option is a zenserp client option that can be passed to `NewClient`
type Option func(*client)

// WithHTTPClient is an option to configure zenserp client with custom http client
func WithHTTPClient(httpClient *http.Client) Option {
	return func(a *client) {
		a.httpClient = httpClient
	}
}

// Client interface
type Client interface {
	// Search provides zenserp search functionality
	Search(ctx context.Context, query string, num int) (*QueryResult, error)
}

type client struct {
	apiKey     string
	baseURL    *url.URL
	httpClient *http.Client
}

// NewClient instantiates a zenserp client
func NewClient(apiKey string, opts ...Option) (Client, error) {
	baseURL, err := url.Parse(zenserpBaseURL)
	if err != nil {
		return nil, fmt.Errorf("error parsing zenser Base URL (%w)", err)
	}

	c := &client{
		apiKey:     apiKey,
		baseURL:    baseURL,
		httpClient: pkgHttp.DefaultHTTPClient(time.Duration(10 * time.Second)),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c, nil
}
