package notion

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mkfsn/notion-go/rest"
	"github.com/stretchr/testify/assert"
)

func Test_blocksChildrenClient_List(t *testing.T) {
	type fields struct {
		restClient      rest.Interface
		mockHTTPHandler http.Handler
	}

	type args struct {
		ctx    context.Context
		params BlocksChildrenListParameters
	}

	type wants struct {
		response *BlocksChildrenListResponse
		err      error
	}

	type test struct {
		name   string
		fields fields
		args   args
		wants  wants
	}

	tests := []test{
		{
			name: "List 3 block children in a page",
			fields: fields{
				restClient: rest.New(),
				mockHTTPHandler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
					assert.Equal(t, http.MethodGet, request.Method)
					assert.Equal(t, "/v1/blocks/b55c9c91-384d-452b-81db-d1ef79372b75/children?page_size=100", request.RequestURI)

					writer.WriteHeader(http.StatusOK)

					_, err := writer.Write([]byte(`{
					    "object": "list",
					    "results": [
					       {
					            "object": "block",
					            "id": "9bc30ad4-9373-46a5-84ab-0a7845ee52e6",
					            "created_time": "2021-03-16T16:31:00.000Z",
					            "last_edited_time": "2021-03-16T16:32:00.000Z",
					            "has_children": false,
					            "type": "heading_2",
					            "heading_2": {
					                "text": [
					                    {
					                        "type": "text",
					                        "text": {
					                            "content": "Lacinato kale",
					                            "link": null
					                        },
					                        "annotations": {
					                            "bold": false,
					                            "italic": false,
					                            "strikethrough": false,
					                            "underline": false,
					                            "code": false,
					                            "color": "default"
					                        },
					                        "plain_text": "Lacinato kale",
					                        "href": null
					                    }
					                ]
					            }
					        },
					        {
					            "object": "block",
					            "id": "7face6fd-3ef4-4b38-b1dc-c5044988eec0",
					            "created_time": "2021-03-16T16:34:00.000Z",
					            "last_edited_time": "2021-03-16T16:36:00.000Z",
					            "has_children": false,
					            "type": "paragraph",
					            "paragraph": {
					                "text": [
					                    {
					                        "type": "text",
					                        "text": {
					                            "content": "Lacinato kale",
					                            "link": {
					                                "url": "https://en.wikipedia.org/wiki/Lacinato_kale"
					                            }
					                        },
					                        "annotations": {
					                            "bold": false,
					                            "italic": false,
					                            "strikethrough": false,
					                            "underline": false,
					                            "code": false,
					                            "color": "default"
					                        },
					                        "plain_text": "Lacinato kale",
					                        "href": "https://en.wikipedia.org/wiki/Lacinato_kale"
					                    },
					                    {
					                        "type": "text",
					                        "text": {
					                            "content": " is a variety of kale with a long tradition in Italian cuisine, especially that of Tuscany. It is also known as Tuscan kale, Italian kale, dinosaur kale, kale, flat back kale, palm tree kale, or black Tuscan palm.",
					                            "link": null
					                        },
					                        "annotations": {
					                            "bold": false,
					                            "italic": false,
					                            "strikethrough": false,
					                            "underline": false,
					                            "code": false,
					                            "color": "default"
					                        },
					                        "plain_text": " is a variety of kale with a long tradition in Italian cuisine, especially that of Tuscany. It is also known as Tuscan kale, Italian kale, dinosaur kale, kale, flat back kale, palm tree kale, or black Tuscan palm.",
					                        "href": null
					                    }
					                ]
					            }
					        },
					        {
					            "object": "block",
					            "id": "7636e2c9-b6c1-4df1-aeae-3ebf0073c5cb",
					            "created_time": "2021-03-16T16:35:00.000Z",
					            "last_edited_time": "2021-03-16T16:36:00.000Z",
					            "has_children": true,
					            "type": "toggle",
					            "toggle": {
					                "text": [
					                    {
					                        "type": "text",
					                        "text": {
					                            "content": "Recipes",
					                            "link": null
					                        },
					                        "annotations": {
					                            "bold": true,
					                            "italic": false,
					                            "strikethrough": false,
					                            "underline": false,
					                            "code": false,
					                            "color": "default"
					                        },
					                        "plain_text": "Recipes",
					                        "href": null
					                    }
					                ]
					            }
					        }
					    ],
					    "next_cursor": null,
					    "has_more": false
					}`))
					assert.NoError(t, err)
				}),
			},
			args: args{
				ctx: context.Background(),
				params: BlocksChildrenListParameters{
					PaginationParameters: PaginationParameters{
						StartCursor: "",
						PageSize:    100,
					},
					BlockID: "b55c9c91-384d-452b-81db-d1ef79372b75",
				},
			},
			wants: wants{
				response: &BlocksChildrenListResponse{
					PaginatedList: PaginatedList{
						Object:     ObjectTypeList,
						HasMore:    false,
						NextCursor: "",
					},
					Results: []Block{
						&Heading2Block{
							BlockBase: BlockBase{
								Object:         ObjectTypeBlock,
								ID:             "9bc30ad4-9373-46a5-84ab-0a7845ee52e6",
								Type:           BlockTypeHeading2,
								CreatedTime:    newTime(time.Date(2021, 3, 16, 16, 31, 0, 0, time.UTC)),
								LastEditedTime: newTime(time.Date(2021, 3, 16, 16, 32, 0, 0, time.UTC)),
								HasChildren:    false,
							},
							Heading2: HeadingBlock{
								Text: []RichText{
									&RichTextText{
										BaseRichText: BaseRichText{
											PlainText: "Lacinato kale",
											Href:      "",
											Type:      RichTextTypeText,
											Annotations: &Annotations{
												Bold:          false,
												Italic:        false,
												Strikethrough: false,
												Underline:     false,
												Code:          false,
												Color:         ColorDefault,
											},
										},
										Text: TextObject{
											Content: "Lacinato kale",
											Link:    nil,
										},
									},
								},
							},
						},
						&ParagraphBlock{
							BlockBase: BlockBase{
								Object:         ObjectTypeBlock,
								ID:             "7face6fd-3ef4-4b38-b1dc-c5044988eec0",
								Type:           BlockTypeParagraph,
								CreatedTime:    newTime(time.Date(2021, 3, 16, 16, 34, 0, 0, time.UTC)),
								LastEditedTime: newTime(time.Date(2021, 3, 16, 16, 36, 0, 0, time.UTC)),
								HasChildren:    false,
							},
							Paragraph: RichTextBlock{
								Text: []RichText{
									&RichTextText{
										BaseRichText: BaseRichText{
											PlainText: "Lacinato kale",
											Href:      "https://en.wikipedia.org/wiki/Lacinato_kale",
											Type:      RichTextTypeText,
											Annotations: &Annotations{
												Bold:          false,
												Italic:        false,
												Strikethrough: false,
												Underline:     false,
												Code:          false,
												Color:         ColorDefault,
											},
										},
										Text: TextObject{
											Content: "Lacinato kale",
											Link: &Link{
												URL: "https://en.wikipedia.org/wiki/Lacinato_kale",
											},
										},
									},

									&RichTextText{
										BaseRichText: BaseRichText{
											PlainText: " is a variety of kale with a long tradition in Italian cuisine, especially that of Tuscany. It is also known as Tuscan kale, Italian kale, dinosaur kale, kale, flat back kale, palm tree kale, or black Tuscan palm.",
											Href:      "",
											Type:      RichTextTypeText,
											Annotations: &Annotations{
												Bold:          false,
												Italic:        false,
												Strikethrough: false,
												Underline:     false,
												Code:          false,
												Color:         ColorDefault,
											},
										},
										Text: TextObject{
											Content: " is a variety of kale with a long tradition in Italian cuisine, especially that of Tuscany. It is also known as Tuscan kale, Italian kale, dinosaur kale, kale, flat back kale, palm tree kale, or black Tuscan palm.",
											Link:    nil,
										},
									},
								},
								Children: []Block{},
							},
						},
						&ToggleBlock{
							BlockBase: BlockBase{
								Object:         ObjectTypeBlock,
								ID:             "7636e2c9-b6c1-4df1-aeae-3ebf0073c5cb",
								Type:           BlockTypeToggle,
								CreatedTime:    newTime(time.Date(2021, 3, 16, 16, 35, 0, 0, time.UTC)),
								LastEditedTime: newTime(time.Date(2021, 3, 16, 16, 36, 0, 0, time.UTC)),
								HasChildren:    true,
							},
							Toggle: RichTextBlock{
								Text: []RichText{
									&RichTextText{
										BaseRichText: BaseRichText{
											PlainText: "Recipes",
											Href:      "",
											Type:      RichTextTypeText,
											Annotations: &Annotations{
												Bold:          true,
												Italic:        false,
												Strikethrough: false,
												Underline:     false,
												Code:          false,
												Color:         ColorDefault,
											},
										},
										Text: TextObject{
											Content: "Recipes",
											Link:    nil,
										},
									},
								},
								Children: []Block{},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockHTTPServer := httptest.NewServer(tt.fields.mockHTTPHandler)
			defer mockHTTPServer.Close()

			d := &blocksChildrenClient{
				restClient: tt.fields.restClient.BaseURL(mockHTTPServer.URL),
			}
			got, err := d.List(tt.args.ctx, tt.args.params)
			if tt.wants.err != nil {
				assert.ErrorIs(t, err, tt.wants.err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.wants.response, got)
		})
	}
}

func Test_blocksChildrenClient_Append(t *testing.T) {
	type fields struct {
		restClient      rest.Interface
		mockHTTPHandler http.Handler
	}

	type args struct {
		ctx    context.Context
		params BlocksChildrenAppendParameters
	}

	type wants struct {
		response *BlocksChildrenAppendResponse
		err      error
	}

	type test struct {
		name   string
		fields fields
		args   args
		wants  wants
	}

	tests := []test{
		{
			name: "Append children successfully",
			fields: fields{
				restClient: rest.New(),
				mockHTTPHandler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
					assert.Equal(t, http.MethodPatch, request.Method)
					assert.Equal(t, "/v1/blocks/9bd15f8d-8082-429b-82db-e6c4ea88413b/children", request.RequestURI)
					assert.Equal(t, "application/json", request.Header.Get("Content-Type"))

					expectedData := `{
						"children": [
							{
								"object": "block",
								"type": "heading_2",
								"heading_2": {
									"text": [{ "type": "text", "text": { "content": "Lacinato kale" } }]
								}
							},
							{
								"object": "block",
								"type": "paragraph",
								"paragraph": {
									"text": [
										{
											"type": "text",
											"text": {
												"content": "Lacinato kale is a variety of kale with a long tradition in Italian cuisine, especially that of Tuscany. It is also known as Tuscan kale, Italian kale, dinosaur kale, kale, flat back kale, palm tree kale, or black Tuscan palm.",
												"link": { "url": "https://en.wikipedia.org/wiki/Lacinato_kale" }
											}
										}
									]
								}
							}
						]
					}`
					b, err := ioutil.ReadAll(request.Body)
					assert.NoError(t, err)
					assert.JSONEq(t, expectedData, string(b))

					writer.WriteHeader(http.StatusOK)

					_, err = writer.Write([]byte(`{
					  "object": "block",
					  "id": "9bd15f8d-8082-429b-82db-e6c4ea88413b",
					  "created_time": "2020-03-17T19:10:04.968Z",
					  "last_edited_time": "2020-03-17T21:49:37.913Z",
					  "has_children": true,
					  "type": "toggle",
					  "toggle": {
					    "text": [
					        {
					            "type": "text",
					            "text": {
					                "content": "Recipes",
					                "link": null
					            },
					            "annotations": {
					                "bold": true,
					                "italic": false,
					                "strikethrough": false,
					                "underline": false,
					                "code": false,
					                "color": "default"
					            },
					            "plain_text": "Recipes",
					            "href": null
					        }
					    ]
					  }
					}`))
					assert.NoError(t, err)
				}),
			},
			args: args{
				ctx: context.Background(),
				params: BlocksChildrenAppendParameters{
					BlockID: "9bd15f8d-8082-429b-82db-e6c4ea88413b",
					Children: []Block{
						&Heading2Block{
							BlockBase: BlockBase{
								Object: ObjectTypeBlock,
								Type:   BlockTypeHeading2,
							},
							Heading2: HeadingBlock{
								Text: []RichText{
									&RichTextText{
										BaseRichText: BaseRichText{
											Type: RichTextTypeText,
										},
										Text: TextObject{
											Content: "Lacinato kale",
										},
									},
								},
							},
						},
						&ParagraphBlock{
							BlockBase: BlockBase{
								Object: ObjectTypeBlock,
								Type:   BlockTypeParagraph,
							},
							Paragraph: RichTextBlock{
								Text: []RichText{
									&RichTextText{
										BaseRichText: BaseRichText{
											Type: RichTextTypeText,
										},
										Text: TextObject{
											Content: "Lacinato kale is a variety of kale with a long tradition in Italian cuisine, especially that of Tuscany. It is also known as Tuscan kale, Italian kale, dinosaur kale, kale, flat back kale, palm tree kale, or black Tuscan palm.",
											Link: &Link{
												URL: "https://en.wikipedia.org/wiki/Lacinato_kale",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			wants: wants{
				response: &BlocksChildrenAppendResponse{
					Block: &ToggleBlock{
						BlockBase: BlockBase{
							Object:         ObjectTypeBlock,
							ID:             "9bd15f8d-8082-429b-82db-e6c4ea88413b",
							Type:           BlockTypeToggle,
							CreatedTime:    newTime(time.Date(2020, 3, 17, 19, 10, 4, 968_000_000, time.UTC)),
							LastEditedTime: newTime(time.Date(2020, 3, 17, 21, 49, 37, 913_000_000, time.UTC)),
							HasChildren:    true,
						},
						Toggle: RichTextBlock{
							Text: []RichText{
								&RichTextText{
									BaseRichText: BaseRichText{
										PlainText: "Recipes",
										Href:      "",
										Type:      RichTextTypeText,
										Annotations: &Annotations{
											Bold:          true,
											Italic:        false,
											Strikethrough: false,
											Underline:     false,
											Code:          false,
											Color:         ColorDefault,
										},
									},
									Text: TextObject{
										Content: "Recipes",
										Link:    nil,
									},
								},
							},
							Children: []Block{},
						},
					},
				},
				err: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockHTTPServer := httptest.NewServer(tt.fields.mockHTTPHandler)
			defer mockHTTPServer.Close()

			d := &blocksChildrenClient{
				restClient: tt.fields.restClient.BaseURL(mockHTTPServer.URL),
			}
			got, err := d.Append(tt.args.ctx, tt.args.params)
			if tt.wants.err != nil {
				assert.ErrorIs(t, err, tt.wants.err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.wants.response, got)
		})
	}
}

func newTime(t time.Time) *time.Time {
	return &t
}
