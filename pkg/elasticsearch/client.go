package elasticsearch

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	es "github.com/elastic/go-elasticsearch/v7"
	"github.com/jponc/rank-app/pkg/awssigner"
)

type Client interface {
	// UpdateMapping updates index mapping
	UpdateMapping(ctx context.Context, index string, body interface{}) error
	// IndexExists checks if ES index currently exists
	IndexExists(ctx context.Context, index string) (*bool, error)
	// CreateIndex creates an ES index
	CreateIndex(ctx context.Context, index string) error
}

type client struct {
	esClient *es.Client
}

// NewClient instantiates ElasticSearch client
func NewClient(esUrl, awsRegion string) (Client, error) {
	address := fmt.Sprintf("https://%s", esUrl)

	// TODO Change to dependency injection
	credentials := credentials.NewEnvCredentials()

	// This uses AWS Signer v4 to sign the request before sending to ES
	esConfig := es.Config{
		Addresses: []string{address},
		Transport: &awssigner.SignV4SDKV1{
			RoundTripper: http.DefaultTransport,
			Credentials:  credentials,
			Region:       awsRegion,
			Now:          time.Now,
			Service:      "es",
		},
	}

	esClient, err := es.NewClient(esConfig)
	if err != nil {
		return nil, fmt.Errorf("cannot create elasticsearch client: %v", err)
	}

	c := &client{
		esClient: esClient,
	}

	return c, nil
}
