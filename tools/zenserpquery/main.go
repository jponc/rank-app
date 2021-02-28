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

func fetchResults(ctx context.Context, client zenserp.Client, keyword string, chResults chan Result, wg *sync.WaitGroup) {
	defer func() {
		fmt.Printf("Finished %s\n", keyword)
		wg.Done()
	}()

	res, err := client.Search(ctx, keyword, 100)
	if err != nil {
		chResults <- Result{
			Keyword: keyword,
			Err:     err,
		}
		return
	}

	chResults <- Result{
		Keyword:       keyword,
		ZenserpResult: res,
	}
	return
}

type Result struct {
	Keyword       string
	ZenserpResult *zenserp.QueryResult
	Err           error
}

func main() {
	startTime := time.Now()
	ctx := context.Background()
	keywords := sampleKeywords
	file := "results.csv"

	var wg sync.WaitGroup

	chResults := make(chan Result, len(keywords))

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

	failedKeywords := make([]string, 0)

	for _, keyword := range keywords {
		wg.Add(1)
		go fetchResults(ctx, client, keyword, chResults, &wg)
	}

	wg.Wait()
	fmt.Println("YO")
	close(chResults)

	records := [][]string{
		{"keyword", "position", "title", "url", "description"},
	}

	for res := range chResults {
		keyword := res.Keyword

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
	}

	w := csv.NewWriter(f)
	err = w.WriteAll(records)
	if err != nil {
		log.Fatal(err)
	}

	endTime := time.Now()
	diff := endTime.Sub(startTime)
	fmt.Printf("\n=> Finished processing top 100 results for %d keywords at %s", len(keywords), diff)
}
