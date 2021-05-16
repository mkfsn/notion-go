package main

import (
	"context"
	"log"
	"os"

	"github.com/mkfsn/notion-go"
)

func main() {
	c := notion.New(notion.WithAuthToken(os.Getenv("NOTION_AUTH_TOKEN")))

	resp, err := c.Databases().Retrieve(context.Background(), notion.DatabasesRetrieveParameters{
		DatabaseID: "def72422-ea36-4c8a-a6f1-a34e11a7fe54",
	})

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%#v\n", resp)
}
