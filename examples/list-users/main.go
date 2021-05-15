package main

import (
	"context"
	"log"
	"os"

	"github.com/mkfsn/notion-go"
)

func main() {
	c := notion.New(notion.WithAuthToken(os.Getenv("NOTION_AUTH_TOKEN")))

	resp, err := c.Users().List(context.Background(), notion.UsersListParameters{
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
