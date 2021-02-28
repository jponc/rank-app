package main

import (
	"context"
	"fmt"

	"github.com/jponc/rank-app/pkg/zenserp"
)

func main() {
	ctx := context.Background()

	config, err := NewConfig()
	if err != nil {
		panic(err)
	}

	client, err := zenserp.NewClient(config.APIKey)
	if err != nil {
		panic(err)
	}

	query := "Water tanks"
	res, err := client.Search(ctx, query)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v", res)
}
