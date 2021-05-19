package main

import (
	"context"
	"log"
	"os"

	"github.com/mkfsn/notion-go"
)

func main() {
	c := notion.New(os.Getenv("NOTION_AUTH_TOKEN"))

	resp, err := c.Blocks().Children().Append(context.Background(),
		notion.BlocksChildrenAppendParameters{
			BlockID: "12e1d803ee234651a125c6ce13ccd58d",
			Children: []notion.Block{
				notion.Heading2Block{
					BlockBase: notion.BlockBase{
						Object: notion.ObjectTypeBlock,
						Type:   notion.BlockTypeHeading2,
					},
					Heading2: notion.HeadingBlock{
						Text: []notion.RichText{
							notion.RichTextText{
								BaseRichText: notion.BaseRichText{
									Type: notion.RichTextTypeText,
								},
								Text: notion.TextObject{
									Content: "Lacinato kale",
								},
							},
						},
					},
				},

				notion.ParagraphBlock{
					BlockBase: notion.BlockBase{
						Object: notion.ObjectTypeBlock,
						Type:   notion.BlockTypeParagraph,
					},
					Paragraph: notion.RichTextBlock{
						Text: []notion.RichText{
							notion.RichTextText{
								BaseRichText: notion.BaseRichText{
									Type: notion.RichTextTypeText,
								},
								Text: notion.TextObject{
									Content: "Lacinato kale is a variety of kale with a long tradition in Italian cuisine, especially that of Tuscany. It is also known as Tuscan kale, Italian kale, dinosaur kale, kale, flat back kale, palm tree kale, or black Tuscan palm.",
									Link: &notion.Link{
										Type: "url",
										URL:  "https://en.wikipedia.org/wiki/Lacinato_kale",
									},
								},
							},
						},
					},
				},
			},
		},
	)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("response: %#v\n", resp)
}
