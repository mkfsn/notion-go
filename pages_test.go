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

func Test_pagesClient_Retrieve(t *testing.T) {
	type fields struct {
		restClient      rest.Interface
		mockHTTPHandler http.Handler
		authToken       string
	}

	type args struct {
		ctx    context.Context
		params PagesRetrieveParameters
	}

	type wants struct {
		response *PagesRetrieveResponse
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
			name: "Retrieve a page by page_id",
			fields: fields{
				restClient: rest.New(),
				authToken:  "6cf01c0d-3b5e-49ec-a45e-43c1879cf41e",
				mockHTTPHandler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
					assert.Equal(t, DefaultNotionVersion, request.Header.Get("Notion-Version"))
					assert.Equal(t, DefaultUserAgent, request.Header.Get("User-Agent"))
					assert.Equal(t, "Bearer 6cf01c0d-3b5e-49ec-a45e-43c1879cf41e", request.Header.Get("Authorization"))

					assert.Equal(t, http.MethodGet, request.Method)
					assert.Equal(t, "/v1/pages/b55c9c91-384d-452b-81db-d1ef79372b75", request.RequestURI)

					writer.WriteHeader(http.StatusOK)

					_, err := writer.Write([]byte(`{
						"object": "page",
						"id": "b55c9c91-384d-452b-81db-d1ef79372b75",
						"created_time": "2020-03-17T19:10:04.968Z",
					  	"last_edited_time": "2020-03-17T21:49:37.913Z",
					  	"properties": {
					    	"Name": {
                          		"id":"title",
								"type":"title",
                          		"title": [
									{
					          			"type": "text", 
										"text": {"content":"Avocado","link":null},
					          			"annotations": {
	          				      			"bold":true,
          				      				"italic":false,
          				      				"strikethrough":false,
          				      				"underline":false,
          				      				"code":false,
          				      				"color":"default"
					          			} 
									}
						  		]
							},
					    	"Description": {
						  		"type":"rich_text",
						  		"rich_text": [
									{
										"type": "text",
										"text": {"content":"Persea americana","link":null},
          				    			"annotations":{
          				      				"bold":false,
          				      				"italic":false,
          				      				"strikethrough":false,
          				      				"underline":false,
          				      				"code":false,
          				      				"color":"default"
          				    			} 
									}
								]
					    	}
					  	}
					}`))
					assert.NoError(t, err)
				}),
			},
			args: args{
				ctx: context.Background(),
				params: PagesRetrieveParameters{
					PageID: "b55c9c91-384d-452b-81db-d1ef79372b75",
				},
			},
			wants: wants{
				response: &PagesRetrieveResponse{
					Page: Page{
						Object: ObjectTypePage,
						ID:     "b55c9c91-384d-452b-81db-d1ef79372b75",
						Parent: nil,
						Properties: map[string]PropertyValue{
							"Name": &TitlePropertyValue{
								basePropertyValue: basePropertyValue{
									ID:   "title",
									Type: PropertyValueTypeTitle,
								},
								Title: []RichText{
									&RichTextText{
										BaseRichText: BaseRichText{
											Href: "",
											Type: RichTextTypeText,
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
											Content: "Avocado",
										},
									},
								},
							},
							"Description": &RichTextPropertyValue{
								basePropertyValue: basePropertyValue{
									Type: PropertyValueTypeRichText,
								},
								RichText: []RichText{
									&RichTextText{
										BaseRichText: BaseRichText{
											Type: RichTextTypeText,
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
											Content: "Persea americana",
											Link:    nil,
										},
									},
								},
							},
						},
						CreatedTime:    time.Date(2020, 3, 17, 19, 10, 4, 968_000_000, time.UTC),
						LastEditedTime: time.Date(2020, 3, 17, 21, 49, 37, 913_000_000, time.UTC),
						Archived:       false,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockHTTPServer := httptest.NewServer(tt.fields.mockHTTPHandler)
			defer mockHTTPServer.Close()

			sut := New(
				tt.fields.authToken,
				WithBaseURL(mockHTTPServer.URL),
			)

			got, err := sut.Pages().Retrieve(tt.args.ctx, tt.args.params)
			if tt.wants.err != nil {
				assert.ErrorIs(t, err, tt.wants.err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.wants.response, got)
		})
	}
}

func Test_pagesClient_Create(t *testing.T) {
	type fields struct {
		restClient      rest.Interface
		mockHTTPHandler http.Handler
		authToken       string
	}

	type args struct {
		ctx    context.Context
		params PagesCreateParameters
	}

	type wants struct {
		response *PagesCreateResponse
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
			name: "Create a new page",
			fields: fields{
				restClient: rest.New(),
				authToken:  "0747d2ee-13f0-47b1-950f-511d2c87180d",
				mockHTTPHandler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
					assert.Equal(t, DefaultNotionVersion, request.Header.Get("Notion-Version"))
					assert.Equal(t, DefaultUserAgent, request.Header.Get("User-Agent"))
					assert.Equal(t, "Bearer 0747d2ee-13f0-47b1-950f-511d2c87180d", request.Header.Get("Authorization"))

					assert.Equal(t, http.MethodPost, request.Method)
					assert.Equal(t, "/v1/pages", request.RequestURI)
					assert.Equal(t, "application/json", request.Header.Get("Content-Type"))

					expectedData := `{
						"parent": { "database_id": "48f8fee9cd794180bc2fec0398253067" },
						"properties": {
							"Name": {
								"title": [
									{
										"text": {
											"content": "Tuscan Kale"
										}
									}
								]
							},
							"Description": {
								"rich_text": [
									{
										"text": {
											"content": "A dark green leafy vegetable"
										}
									}
								]
							},
							"Food group": {
								"select": {
									"name": "Vegetable"
								}
							},
							"Price": { "number": 2.5 }
						},
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
						"object": "page",
					  	"id": "251d2b5f-268c-4de2-afe9-c71ff92ca95c",
						"created_time": "2020-03-17T19:10:04.968Z",
						"last_edited_time": "2020-03-17T21:49:37.913Z",
					  	"parent": {
					    	"type": "database_id",
					    	"database_id": "48f8fee9-cd79-4180-bc2f-ec0398253067"
					  	},
					  	"archived": false,
					  	"properties": {
					    	"Recipes": {
					      		"id": "Ai` + "`" + `L",
								"type": "relation",
								"relation": []
							},
							"Cost of next trip": {
								"id": "R}wl",
								"type": "formula",
								"formula": {
									"type": "number",
									"number": null
								}
							},
							"Photos": {
									"id": "d:Cb",
									"type": "files",
									"files": []
							},
							"Store availability": {
								"id": "jrFQ",
								"type": "multi_select",
								"multi_select": []
							},
							"+1": {
								"id": "k?CE",
								"type": "people",
								"people": []
							},
							"Description": {
								"id": "rT{n",
								"type": "rich_text",
								"rich_text": []
							},
							"In stock": {
								"id": "{>U;",
								"type": "checkbox",
								"checkbox": false
							},
							"Name": {
								"id": "title",
								"type": "title",
								"title": [
									{
										"type": "text",
										"text": {
											"content": "Tuscan Kale",
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
										"plain_text": "Tuscan Kale",
										"href": null
									}
								]
							}
						}
					}`))
					assert.NoError(t, err)
				}),
			},
			args: args{
				ctx: context.Background(),
				params: PagesCreateParameters{
					Parent: &DatabaseParentInput{
						DatabaseID: "48f8fee9cd794180bc2fec0398253067",
					},
					Properties: map[string]PropertyValue{
						"Name": &TitlePropertyValue{
							Title: []RichText{
								&RichTextText{
									Text: TextObject{
										Content: "Tuscan Kale",
									},
								},
							},
						},
						"Description": &RichTextPropertyValue{
							basePropertyValue: basePropertyValue{},
							RichText: []RichText{
								&RichTextText{
									BaseRichText: BaseRichText{},
									Text: TextObject{
										Content: "A dark green leafy vegetable",
									},
								},
							},
						},
						"Food group": &SelectPropertyValue{
							basePropertyValue: basePropertyValue{},
							Select: SelectPropertyValueOption{
								Name: "Vegetable",
							},
						},
						"Price": &NumberPropertyValue{
							Number: 2.5,
						},
					},
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
				response: &PagesCreateResponse{
					Page: Page{
						Object: ObjectTypePage,
						ID:     "251d2b5f-268c-4de2-afe9-c71ff92ca95c",
						Parent: &DatabaseParent{
							baseParent: baseParent{
								Type: ParentTypeDatabase,
							},
							DatabaseID: "48f8fee9-cd79-4180-bc2f-ec0398253067",
						},
						Properties: map[string]PropertyValue{
							"Recipes": &RelationPropertyValue{
								basePropertyValue: basePropertyValue{
									ID:   "Ai`L",
									Type: PropertyValueTypeRelation,
								},
								Relation: []PageReference{},
							},
							"Cost of next trip": &FormulaPropertyValue{
								basePropertyValue: basePropertyValue{
									ID:   "R}wl",
									Type: PropertyValueTypeFormula,
								},
								Formula: &NumberFormulaValue{
									baseFormulaValue: baseFormulaValue{
										Type: FormulaValueTypeNumber,
									},
									Number: nil,
								},
							},
							"Photos": &FilesPropertyValue{
								basePropertyValue: basePropertyValue{
									ID:   "d:Cb",
									Type: PropertyValueTypeFiles,
								},
								Files: []File{},
							},
							"Store availability": &MultiSelectPropertyValue{
								basePropertyValue: basePropertyValue{
									ID:   "jrFQ",
									Type: PropertyValueTypeMultiSelect,
								},
								MultiSelect: []MultiSelectPropertyValueOption{},
							},
							"+1": &PeoplePropertyValue{
								basePropertyValue: basePropertyValue{
									ID:   "k?CE",
									Type: PropertyValueTypePeople,
								},
								People: []User{},
							},
							"Description": &RichTextPropertyValue{
								basePropertyValue: basePropertyValue{
									ID:   "rT{n",
									Type: PropertyValueTypeRichText,
								},
								RichText: []RichText{},
							},
							"In stock": &CheckboxPropertyValue{
								basePropertyValue: basePropertyValue{
									ID:   "{>U;",
									Type: PropertyValueTypeCheckbox,
								},
								Checkbox: false,
							},
							"Name": &TitlePropertyValue{
								basePropertyValue: basePropertyValue{
									ID:   "title",
									Type: PropertyValueTypeTitle,
								},
								Title: []RichText{
									&RichTextText{
										BaseRichText: BaseRichText{
											PlainText: "Tuscan Kale",
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
											Content: "Tuscan Kale",
											Link:    nil,
										},
									},
								},
							},
						},
						CreatedTime:    time.Date(2020, 3, 17, 19, 10, 04, 968_000_000, time.UTC),
						LastEditedTime: time.Date(2020, 3, 17, 21, 49, 37, 913_000_000, time.UTC),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockHTTPServer := httptest.NewServer(tt.fields.mockHTTPHandler)
			defer mockHTTPServer.Close()

			sut := New(
				tt.fields.authToken,
				WithBaseURL(mockHTTPServer.URL),
			)

			got, err := sut.Pages().Create(tt.args.ctx, tt.args.params)
			if tt.wants.err != nil {
				assert.ErrorIs(t, err, tt.wants.err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.wants.response, got)
		})
	}
}

func Test_pagesClient_Update(t *testing.T) {
	type fields struct {
		restClient      rest.Interface
		mockHTTPHandler http.Handler
		authToken       string
	}

	type args struct {
		ctx    context.Context
		params PagesUpdateParameters
	}

	type wants struct {
		response *PagesUpdateResponse
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
			name: "Update a page properties",
			fields: fields{
				restClient: rest.New(),
				authToken:  "4ad4d7a9-8b66-4dda-b9a1-2bc98134ee14",
				mockHTTPHandler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
					assert.Equal(t, DefaultNotionVersion, request.Header.Get("Notion-Version"))
					assert.Equal(t, DefaultUserAgent, request.Header.Get("User-Agent"))
					assert.Equal(t, "Bearer 4ad4d7a9-8b66-4dda-b9a1-2bc98134ee14", request.Header.Get("Authorization"))

					assert.Equal(t, http.MethodPatch, request.Method)
					assert.Equal(t, "/v1/pages/60bdc8bd-3880-44b8-a9cd-8a145b3ffbd7", request.RequestURI)
					assert.Equal(t, "application/json", request.Header.Get("Content-Type"))

					expectedData := `{
					  "properties": {
					    "In stock": { "checkbox": true }
					  }
					}`
					b, err := ioutil.ReadAll(request.Body)
					assert.NoError(t, err)
					assert.JSONEq(t, expectedData, string(b))

					writer.WriteHeader(http.StatusOK)

					_, err = writer.Write([]byte(`{
						  "object": "page",
						  "id": "60bdc8bd-3880-44b8-a9cd-8a145b3ffbd7",
							"created_time": "2020-03-17T19:10:04.968Z",
							"last_edited_time": "2020-03-17T21:49:37.913Z",
						  "parent": {
						    "type": "database_id",
						    "database_id": "48f8fee9-cd79-4180-bc2f-ec0398253067"
						  },
						  "archived": false,
						  "properties": {
						    "In stock": {
						      "id": "{>U;",
						      "type": "checkbox",
						      "checkbox": true
						    },
						    "Name": {
						      "id": "title",
						      "type": "title",
						      "title": [
						        {
						          "type": "text",
						          "text": {
						            "content": "Avocado",
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
						          "plain_text": "Avocado",
						          "href": null
						        }
						      ]
						    }
						  }
						}`))
					assert.NoError(t, err)
				}),
			},
			args: args{
				ctx: context.Background(),
				params: PagesUpdateParameters{
					PageID: "60bdc8bd-3880-44b8-a9cd-8a145b3ffbd7",
					Properties: map[string]PropertyValue{
						"In stock": &CheckboxPropertyValue{
							Checkbox: true,
						},
					},
				},
			},
			wants: wants{
				response: &PagesUpdateResponse{
					Page: Page{
						Object: ObjectTypePage,
						ID:     "60bdc8bd-3880-44b8-a9cd-8a145b3ffbd7",
						Parent: &DatabaseParent{
							baseParent: baseParent{
								Type: ParentTypeDatabase,
							},
							DatabaseID: "48f8fee9-cd79-4180-bc2f-ec0398253067",
						},
						Properties: map[string]PropertyValue{
							"In stock": &CheckboxPropertyValue{
								basePropertyValue: basePropertyValue{
									ID:   "{>U;",
									Type: PropertyValueTypeCheckbox,
								},
								Checkbox: true,
							},
							"Name": &TitlePropertyValue{
								basePropertyValue: basePropertyValue{
									ID:   "title",
									Type: PropertyValueTypeTitle,
								},
								Title: []RichText{
									&RichTextText{
										BaseRichText: BaseRichText{
											PlainText: "Avocado",
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
											Content: "Avocado",
											Link:    nil,
										},
									},
								},
							},
						},
						CreatedTime:    time.Date(2020, 3, 17, 19, 10, 4, 968_000_000, time.UTC),
						LastEditedTime: time.Date(2020, 3, 17, 21, 49, 37, 913_000_000, time.UTC),
						Archived:       false,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockHTTPServer := httptest.NewServer(tt.fields.mockHTTPHandler)
			defer mockHTTPServer.Close()

			sut := New(
				tt.fields.authToken,
				WithBaseURL(mockHTTPServer.URL),
			)

			got, err := sut.Pages().Update(tt.args.ctx, tt.args.params)
			if tt.wants.err != nil {
				assert.ErrorIs(t, err, tt.wants.err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.wants.response, got)
		})
	}
}
