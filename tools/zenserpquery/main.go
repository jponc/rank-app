package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	pkgHttp "github.com/jponc/rank-app/pkg/http"
	"github.com/jponc/rank-app/pkg/zenserp"
	log "github.com/sirupsen/logrus"
)

func fetchResults(ctx context.Context, client zenserp.Client, keyword string, chResults chan result) {
	defer func() {
		fmt.Printf("Finished %s\n", keyword)
	}()

	res, err := client.Search(ctx, keyword, 100)
	if err != nil {
		chResults <- result{
			Keyword: keyword,
			Err:     err,
		}
		return
	}

	chResults <- result{
		Keyword:       keyword,
		ZenserpResult: res,
	}
	return
}

type result struct {
	Keyword       string
	ZenserpResult *zenserp.QueryResult
	Err           error
}

func main() {
	startTime := time.Now()
	ctx := context.Background()

	var wg sync.WaitGroup

	keywords := []string{
		"sydney",
		"craft beers",
		"gin tonic",
		"hello world",
	}

	file := "results.csv"

	workerCount := 2
	chResults := make(chan result, workerCount)

	f, err := os.Create(file)
	defer f.Close()

	if err != nil {
		log.Fatalln("failed to open file", err)
	}

	config, err := NewConfig()
	if err != nil {
		log.Fatalln(err)
	}

	httpClient := pkgHttp.DefaultHTTPClient(time.Duration(1 * time.Minute))
	client, err := zenserp.NewClient(config.APIKey, zenserp.WithHTTPClient(httpClient))
	if err != nil {
		log.Fatalln(err)
	}

	finishAllWaitGroups := func() {
		for workerIdx := 0; workerIdx < workerCount; workerIdx++ {
			wg.Done()
		}
	}

	records := [][]string{
		{"keyword", "position", "title", "url", "description"},
	}

	failedKeywords := make([]string, 0)
	done := 0

	for workerIdx := 0; workerIdx < workerCount; workerIdx++ {
		wg.Add(1)
		go func(workerIdx int) {
			for res := range chResults {
				keyword := res.Keyword
				fmt.Printf("Worker %d processed %s", workerIdx, keyword)

				if res.Err != nil {
					failedKeywords = append(failedKeywords, keyword)
				} else {
					zRes := res.ZenserpResult
					for _, resultItem := range zRes.ResulItems {
						records = append(records, []string{
							keyword,
							strconv.Itoa(resultItem.Position),
							resultItem.Title,
							resultItem.URL,
							resultItem.Description,
						})
					}
				}

				done++
				if done == len(keywords) {
					finishAllWaitGroups()
				}
			}
		}(workerIdx)
	}

	for _, keyword := range keywords {
		go fetchResults(ctx, client, keyword, chResults)
	}

	wg.Wait()
	close(chResults)

	w := csv.NewWriter(f)
	err = w.WriteAll(records)
	if err != nil {
		log.Fatal(err)
	}

	endTime := time.Now()
	diff := endTime.Sub(startTime)
	fmt.Printf("\n=> Finished processing top 100 results for %d keywords at %s", len(keywords), diff)
}
