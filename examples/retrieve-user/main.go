package main

import (
	"context"
	"log"
	"os"

	"github.com/mkfsn/notion-go"
)

func main() {
	c := notion.New(notion.WithAuthToken(os.Getenv("NOTION_AUTH_TOKEN")))

	resp, err := c.Users().Retrieve(context.Background(), notion.UsersRetrieveParameters{UserID: "8cd69bf3-1532-43d2-9b11-9803c813d607"})

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%#v\n", resp)
}
