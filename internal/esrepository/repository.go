package esrepository

import (
	"context"

	"github.com/jponc/rank-app/pkg/elasticsearch"
)

const (
	resultItemIndexName = "result-item"
)

// Repository interface
type Repository interface {
	// UpdateResultItemIndexMapping updates ES mapping
	UpdateResultItemIndexMapping(ctx context.Context) error
}

type repository struct {
	esClient elasticsearch.Client
}

func NewRepository(esClient elasticsearch.Client) (Repository, error) {
	// NewClient instantiates a repository
	r := &repository{
		esClient: esClient,
	}

	return r, nil
}
