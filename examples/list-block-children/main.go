package main

import (
	"context"
	"log"
	"os"

	"github.com/mkfsn/notion-go"
)

func main() {
	c := notion.New(os.Getenv("NOTION_AUTH_TOKEN"))

	resp, err := c.Blocks().Children().List(context.Background(), notion.BlocksChildrenListParameters{
		BlockID: "12e1d803ee234651a125c6ce13ccd58d"},
	)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("response: %#v\n", resp)
}
