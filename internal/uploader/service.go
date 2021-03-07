package uploader

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/jponc/rank-app/internal/repository"
	"github.com/jponc/rank-app/internal/s3repository"
)

type Service interface {
	UploadLatestCSV(ctx context.Context)
}

type service struct {
	repository   repository.Repository
	s3repository s3repository.Repository
}

// NewService instantiates a new service
func NewService(s3repository s3repository.Repository, repository repository.Repository) Service {
	return &service{
		s3repository: s3repository,
		repository:   repository,
	}
}

func (s *service) UploadLatestCSV(ctx context.Context) {
	if s.repository == nil {
		log.Fatal("repository not defined")
	}

	if s.s3repository == nil {
		log.Fatal("s3repository not defined")
	}

	latestCsv, err := s.getLatestCsv()
	if err != nil {
		log.Fatal(err)
	}

	err = s.s3repository.UploadLatestCSV(latestCsv)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *service) getLatestCsv() (string, error) {
	crawlResults, err := s.repository.GetLatestCrawlResults()
	if err != nil {
		return "", err
	}

	pairs := [][]string{
		{
			"query",
			"search_engine",
			"device",
			"item_position",
			"item_title",
			"item_url",
			"item_description",
			"result_url",
			"result_id",
			"created_at",
		},
	}

	for _, crawlResult := range *crawlResults {
		for _, item := range crawlResult.Items {
			row := []string{
				crawlResult.Query,
				crawlResult.SearchEngine,
				crawlResult.Device,
				fmt.Sprint(item.Position),
				item.Title,
				item.URL,
				item.Description,
				crawlResult.URL,
				crawlResult.ID.String(),
				crawlResult.CreatedAt.String(),
			}

			pairs = append(pairs, row)
		}
	}

	b := new(bytes.Buffer)
	w := csv.NewWriter(b)
	err = w.WriteAll(pairs)

	if err != nil {
		return "", err
	}

	csvString := b.String()

	return csvString, nil
}
