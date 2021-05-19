package main

import (
	"context"
	"log"
	"os"

	"github.com/mkfsn/notion-go"
)

func main() {
	c := notion.New(os.Getenv("NOTION_AUTH_TOKEN"))

	page, err := c.Pages().Retrieve(context.Background(), notion.PagesRetrieveParameters{
		PageID: "676aa7b7-2bba-4b5b-9fd6-b43f5543482d"},
	)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("page: %#v\n", page)
}
