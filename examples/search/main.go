package main

import (
	"context"
	"log"
	"os"

	"github.com/mkfsn/notion-go"
	"github.com/mkfsn/notion-go/typed"
)

func main() {
	c := notion.New(os.Getenv("NOTION_AUTH_TOKEN"))

	resp, err := c.Search(context.Background(), notion.SearchParameters{
		Query: "フィリスのアトリエ",
		Sort: notion.SearchSort{
			Direction: typed.SearchSortDirectionAscending,
			Timestamp: typed.SearchSortTimestampLastEditedTime,
		},
		Filter: notion.SearchFilter{
			Property: typed.SearchFilterPropertyObject,
			Value:    typed.SearchFilterValuePage,
		},
	})
	if err != nil {
		log.Fatalf("error: %s\n", err)
	}

	log.Printf("response: %#v\n", resp)
}
