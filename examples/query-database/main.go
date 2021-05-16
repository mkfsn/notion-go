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

	keyword := "medium.com"

	resp, err := c.Databases().Query(context.Background(), notion.DatabasesQueryParameters{
		DatabaseID: "def72422-ea36-4c8a-a6f1-a34e11a7fe54",
		Filter: notion.CompoundFilter{
			Or: []notion.Filter{
				notion.SingleTextFilter{
					SinglePropertyFilter: notion.SinglePropertyFilter{
						Property: "URL",
					},
					URL: &notion.TextFilter{
						Contains: &keyword,
					},
				},
			},
		},
		Sorts: []notion.Sort{
			{
				Property:  "Created",
				Direction: typed.SortDirectionAscending,
			},
		},
	})

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%#v\n", resp)
}
