package main

import (
	"context"
	"log"
	"os"

	"github.com/mkfsn/notion-go"
)

func main() {
	c := notion.New(os.Getenv("NOTION_AUTH_TOKEN"))

	resp, err := c.Search(context.Background(), notion.SearchParameters{
		Query: "フィリスのアトリエ",
		Sort: notion.SearchSort{
			Direction: notion.SearchSortDirectionAscending,
			Timestamp: notion.SearchSortTimestampLastEditedTime,
		},
		Filter: notion.SearchFilter{
			Property: notion.SearchFilterPropertyObject,
			Value:    notion.SearchFilterValuePage,
		},
	})
	if err != nil {
		log.Fatalf("error: %s\n", err)
	}

	for _, object := range resp.Results {
		log.Printf("object: %#v\n", object)
	}
}
