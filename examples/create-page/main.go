package main

import (
	"context"
	"log"
	"os"

	"github.com/mkfsn/notion-go"
)

func main() {
	c := notion.New(notion.WithAuthToken(os.Getenv("NOTION_AUTH_TOKEN")))

	page, err := c.Pages().Create(context.Background(),
		notion.PagesCreateParameters{
			Parent: notion.DatabaseParentInput{
				DatabaseID: "aee104a17e554846bea3536712bfca2c",
			},

			Properties: map[string]notion.PropertyValue{
				"Name": notion.TitlePropertyValue{
					Title: []notion.RichText{
						notion.RichTextText{Text: notion.TextObject{Content: "Tuscan Kale"}},
					},
				},

				"Description": notion.RichTextPropertyValue{
					RichText: []notion.RichText{
						notion.RichTextText{Text: notion.TextObject{Content: " dark green leafy vegetable"}},
					},
				},

				"Food group": notion.SelectPropertyValue{
					Select: notion.SelectPropertyValueOption{
						Name: "Vegetable",
					},
				},

				"Price": notion.NumberPropertyValue{
					Number: 2.5,
				},
			},

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

	log.Printf("page: %#v\n", page)
}
