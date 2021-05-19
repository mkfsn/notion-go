package main

import (
	"context"
	"log"
	"os"

	"github.com/mkfsn/notion-go"
)

func main() {
	c := notion.New(os.Getenv("NOTION_AUTH_TOKEN"))

	resp, err := c.Databases().List(context.Background(), notion.DatabasesListParameters{
		PaginationParameters: notion.PaginationParameters{
			StartCursor: "",
			PageSize:    1,
		},
	})

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%#v\n", resp)
}
