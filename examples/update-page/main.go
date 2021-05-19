package main

import (
	"context"
	"log"
	"os"

	"github.com/mkfsn/notion-go"
)

func main() {
	c := notion.New(os.Getenv("NOTION_AUTH_TOKEN"))

	page, err := c.Pages().Update(context.Background(),
		notion.PagesUpdateParameters{
			PageID: "6eaac3811afd4f368209b572e13eace4",
			Properties: map[string]notion.PropertyValue{
				"In stock": notion.CheckboxPropertyValue{
					Checkbox: true,
				},

				"Price": notion.NumberPropertyValue{
					Number: 30,
				},
			},
		},
	)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("page: %#v\n", page)
}
