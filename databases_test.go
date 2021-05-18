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

func Test_databasesClient_List(t *testing.T) {
	type fields struct {
		restClient      rest.Interface
		mockHTTPHandler http.Handler
	}

	type args struct {
		ctx    context.Context
		params DatabasesListParameters
	}

	type wants struct {
		response *DatabasesListResponse
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
			name: "List two databases in one page",
			fields: fields{
				restClient: rest.New(),
				mockHTTPHandler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
					assert.Equal(t, http.MethodGet, request.Method)
					assert.Equal(t, "/v1/databases?page_size=2", request.RequestURI)

					writer.WriteHeader(http.StatusOK)

					_, err := writer.Write([]byte(`{
					    "results": [
					        {
					            "object": "database",
					            "id": "668d797c-76fa-4934-9b05-ad288df2d136",
					            "properties": {
					                "Name": {
					                    "type": "title",
					                    "title": {}
					                },
					                "Description": {
					                    "type": "rich_text",
					                    "rich_text": {}
					                }
					            }
					        },
					        {
					            "object": "database",
					            "id": "74ba0cb2-732c-4d2f-954a-fcaa0d93a898",
					            "properties": {
					                "Name": {
					                    "type": "title",
					                    "title": {}
					                },
					                "Description": {
					                    "type": "rich_text",
					                    "rich_text": {}
					                }
					            }
					        }
					    ],
					    "next_cursor": "MTY3NDE4NGYtZTdiYy00NzFlLWE0NjctODcxOTIyYWU3ZmM3",
					    "has_more": false
					}`))
					assert.NoError(t, err)
				}),
			},
			args: args{
				ctx: context.Background(),
				params: DatabasesListParameters{
					PaginationParameters: PaginationParameters{
						StartCursor: "",
						PageSize:    2,
					},
				},
			},
			wants: wants{
				response: &DatabasesListResponse{
					PaginatedList: PaginatedList{
						// FIXME: This should be ObjectTypeList but the example does not provide the object key-value
						Object:     "",
						HasMore:    false,
						NextCursor: "MTY3NDE4NGYtZTdiYy00NzFlLWE0NjctODcxOTIyYWU3ZmM3",
					},
					Results: []Database{
						{
							Object: ObjectTypeDatabase,
							ID:     "668d797c-76fa-4934-9b05-ad288df2d136",
							// FIXME: The example seems to have a invalid title thus I remove it for now, but need to check
							// what is the expected title
							Title: []RichText{},
							Properties: map[string]Property{
								"Name": &TitleProperty{
									baseProperty: baseProperty{
										ID:   "",
										Type: PropertyTypeTitle,
									},
									Title: map[string]interface{}{},
								},
								// FIXME: The example seems to have a invalid type of the description thus I change it
								// to `rich_text` for now, but need to check what is the expected type
								"Description": &RichTextProperty{
									baseProperty: baseProperty{
										ID:   "",
										Type: PropertyTypeRichText,
									},
									RichText: map[string]interface{}{},
								},
							},
						},
						{
							Object: ObjectTypeDatabase,
							ID:     "74ba0cb2-732c-4d2f-954a-fcaa0d93a898",
							// FIXME: The example seems to have a invalid title thus I remove it for now, but need to check
							// what is the expected title
							Title: []RichText{},
							Properties: map[string]Property{
								"Name": &TitleProperty{
									baseProperty: baseProperty{
										ID:   "",
										Type: PropertyTypeTitle,
									},
									Title: map[string]interface{}{},
								},
								"Description": &RichTextProperty{
									baseProperty: baseProperty{
										ID:   "",
										Type: PropertyTypeRichText,
									},
									RichText: map[string]interface{}{},
								},
							},
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

			d := &databasesClient{
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

func Test_databasesClient_Query(t *testing.T) {
	type fields struct {
		restClient      rest.Interface
		mockHTTPHandler http.Handler
	}

	type args struct {
		ctx    context.Context
		params DatabasesQueryParameters
	}

	type wants struct {
		response *DatabasesQueryResponse
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
			name: "Query database with filter and sort",
			fields: fields{
				restClient: rest.New(),
				mockHTTPHandler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
					assert.Equal(t, http.MethodPost, request.Method)
					assert.Equal(t, "/v1/databases/897e5a76ae524b489fdfe71f5945d1af/query", request.RequestURI)

					expectedData := `{
						"filter": {
						    "or": [
						      	{
						        	"property": "In stock",
							  		"checkbox": {
							  			"equals": true
							  		}
						      	},
						      	{
							  		"property": "Cost of next trip",
							  		"number": {
							  			"greater_than_or_equal_to": 2
							  		}
							  	}
							]
						},
					 	"sorts": [
					    	{
					      		"property": "Last ordered",
					      		"direction": "ascending"
					    	}
					  	]
					}`

					b, err := ioutil.ReadAll(request.Body)
					assert.NoError(t, err)
					assert.JSONEq(t, expectedData, string(b))

					writer.WriteHeader(http.StatusOK)

					_, err = writer.Write([]byte(`{
						"object": "list",
						"results": [
						  {
						    "object": "page",
						    "id": "2e01e904-febd-43a0-ad02-8eedb903a82c",
						    "created_time": "2020-03-17T19:10:04.968Z",
						    "last_edited_time": "2020-03-17T21:49:37.913Z",
						    "parent": {
						      "type": "database_id",
						      "database_id": "897e5a76-ae52-4b48-9fdf-e71f5945d1af"
						    },
						    "archived": false,
						    "properties": {
						      "Recipes": {
						        "id": "Ai` + "`" + `L",
											"type": "relation",
												"relation": [
									{
										"id": "796659b4-a5d9-4c64-a539-06ac5292779e"
									},
									{
										"id": "79e63318-f85a-4909-aceb-96a724d1021c"
									}
									]
									},
									"Cost of next trip": {
										"id": "R}wl",
										"type": "formula",
										"formula": {
											"type": "number",
											"number": 2
										}
									},
									"Last ordered": {
										"id": "UsKi",
										"type": "date",
										"date": {
											"start": "2020-10-07",
											"end": null
										}
									},
									"In stock": {
										"id": "{>U;",
										"type": "checkbox",
										"checkbox": false
									}
								}
							}
						],
						"has_more": false,
						"next_cursor": null
					}`))
					assert.NoError(t, err)
				}),
			},
			args: args{
				ctx: context.Background(),
				params: DatabasesQueryParameters{
					PaginationParameters: PaginationParameters{},
					DatabaseID:           "897e5a76ae524b489fdfe71f5945d1af",
					Filter: CompoundFilter{
						Or: []Filter{
							&SingleCheckboxFilter{
								SinglePropertyFilter: SinglePropertyFilter{
									Property: "In stock",
								},
								Checkbox: CheckboxFilter{
									Equals: true,
								},
							},
							&SingleNumberFilter{
								SinglePropertyFilter: SinglePropertyFilter{
									Property: "Cost of next trip",
								},
								Number: NumberFilter{
									GreaterThanOrEqualTo: newFloat64(2),
								},
							},
						},
					},
					Sorts: []Sort{
						{
							Property:  "Last ordered",
							Direction: SortDirectionAscending,
						},
					},
				},
			},
			wants: wants{
				response: &DatabasesQueryResponse{
					PaginatedList: PaginatedList{
						Object:     ObjectTypeList,
						HasMore:    false,
						NextCursor: "",
					},
					Results: []Page{
						{
							Object: ObjectTypePage,
							ID:     "2e01e904-febd-43a0-ad02-8eedb903a82c",
							Parent: &DatabaseParent{
								baseParent: baseParent{
									Type: ParentTypeDatabase,
								},
								DatabaseID: "897e5a76-ae52-4b48-9fdf-e71f5945d1af",
							},
							Properties: map[string]PropertyValue{
								"Recipes": &RelationPropertyValue{
									basePropertyValue: basePropertyValue{
										ID:   "Ai`L",
										Type: PropertyValueTypeRelation,
									},
									Relation: []PageReference{
										{ID: "796659b4-a5d9-4c64-a539-06ac5292779e"},
										{ID: "79e63318-f85a-4909-aceb-96a724d1021c"},
									},
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
										Number: newFloat64(2),
									},
								},
								"Last ordered": &DatePropertyValue{
									basePropertyValue: basePropertyValue{
										ID:   "UsKi",
										Type: PropertyValueTypeDate,
									},
									Date: Date{
										Start: "2020-10-07",
										End:   nil,
									},
								},
								"In stock": &CheckboxPropertyValue{
									basePropertyValue: basePropertyValue{
										ID:   "{>U;",
										Type: PropertyValueTypeCheckbox,
									},
									Checkbox: false,
								},
							},
							CreatedTime:    time.Date(2020, 3, 17, 19, 10, 4, 968_000_000, time.UTC),
							LastEditedTime: time.Date(2020, 3, 17, 21, 49, 37, 913_000_000, time.UTC),
							Archived:       false,
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockHTTPServer := httptest.NewServer(tt.fields.mockHTTPHandler)

			d := &databasesClient{
				restClient: tt.fields.restClient.BaseURL(mockHTTPServer.URL),
			}
			got, err := d.Query(tt.args.ctx, tt.args.params)
			if tt.wants.err != nil {
				assert.ErrorIs(t, err, tt.wants.err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.wants.response, got)
		})
	}
}

func Test_databasesClient_Retrieve(t *testing.T) {
	type fields struct {
		restClient      rest.Interface
		mockHTTPHandler http.Handler
	}

	type args struct {
		ctx    context.Context
		params DatabasesRetrieveParameters
	}

	type wants struct {
		response *DatabasesRetrieveResponse
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
			name: "Retrieve database",
			fields: fields{
				restClient: rest.New(),
				mockHTTPHandler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
					assert.Equal(t, http.MethodGet, request.Method)
					assert.Equal(t, "/v1/databases/668d797c-76fa-4934-9b05-ad288df2d136", request.RequestURI)

					writer.WriteHeader(http.StatusOK)

					_, err := writer.Write([]byte(`{
						"object": "database",
						"id": "668d797c-76fa-4934-9b05-ad288df2d136",
						"created_time": "2020-03-17T19:10:04.968Z",
						"last_edited_time": "2020-03-17T21:49:37.913Z",
						"title": [
							{
							    "type": "text",
							    "text": {
							        "content": "Grocery List",
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
							    "plain_text": "Grocery List",
							    "href": null
							}
						],
						"properties": {
						    "Name": {
						      "id": "title",
						      "type": "title",
						      "title": {}
						    },
						    "Description": {
						      "id": "J@cS",
						      "type": "rich_text",
						      "rich_text": {}
						    },
						    "In stock": {
						     	"id": "{xY` + "`" + `",
								"type": "checkbox",
								"checkbox": {}
							},
							"Food group": {
								"id": "TJmr",
								"type": "select",
								"select": {
									"options": [
										{
											"id": "96eb622f-4b88-4283-919d-ece2fbed3841",
											"name": "🥦Vegetable",
											"color": "green"
										},
										{
											"id": "bb443819-81dc-46fb-882d-ebee6e22c432",
											"name": "🍎Fruit",
											"color": "red"
										},
										{
											"id": "7da9d1b9-8685-472e-9da3-3af57bdb221e",
											"name": "💪Protein",
											"color": "yellow"
										}
									]
								}
							},
							"Price": {
								"id": "cU^N",
								"type": "number",
								"number": {
									"format": "dollar"
								}
							},
							"Cost of next trip": {
								"id": "p:sC",
								"type": "formula",
								"formula": {
									"expression": "if(prop(\"In stock\"), 0, prop(\"Price\"))"
								}
							},
							"Last ordered": {
								"id": "]\\R[",
								"type": "date",
								"date": {}
							},
							"Meals": {
								"type": "relation",
								"relation": {
									"database_id": "668d797c-76fa-4934-9b05-ad288df2d136",
									"synced_property_name": null
								}
							},
							"Number of meals": {
								"id": "Z\\Eh",
								"type": "rollup",
								"rollup": {
									"rollup_property_name": "Name",
									"relation_property_name": "Meals",
									"rollup_property_id": "title",
									"relation_property_id": "mxp^",
									"function": "count"
								}
							},
							"Store availability": {
								"type": "multi_select",
								"multi_select": {
									"options": [
											{
												"id": "d209b920-212c-4040-9d4a-bdf349dd8b2a",
												"name": "Duc Loi Market",
												"color": "blue"
											},
											{
												"id": "70104074-0f91-467b-9787-00d59e6e1e41",
												"name": "Rainbow Grocery",
												"color": "gray"
											},
											{
												"id": "e6fd4f04-894d-4fa7-8d8b-e92d08ebb604",
												"name": "Nijiya Market",
												"color": "purple"
											},
											{
												"id": "6c3867c5-d542-4f84-b6e9-a420c43094e7",
												"name": "Gus's Community Market",
												"color": "yellow"
											}
									]
								}
							},
							"+1": {
								"id": "aGut",
								"type": "people",
								"people": {}
							},
							"Photo": {
								"id": "aTIT",
								"type": "file",
								"file": {}
							}
						}
					}`))
					// FIXME: In the API reference https://developers.notion.com/reference/get-database
					// the options in the multi_select property is an array of array which does not make sense
					// to me, thus changing it to an array but need to verify this.

					// FIXME: In the API reference https://developers.notion.com/reference/get-database
					// the Photo Property has `files` type but in the documentation,
					// https://developers.notion.com/reference/database, the type should be `file`,
					// Following the documentation, the JSON string above has been changed to type `file`.
					assert.NoError(t, err)
				}),
			},
			args: args{
				ctx: context.Background(),
				params: DatabasesRetrieveParameters{
					DatabaseID: "668d797c-76fa-4934-9b05-ad288df2d136",
				},
			},
			wants: wants{
				response: &DatabasesRetrieveResponse{
					Database: Database{
						Object:         ObjectTypeDatabase,
						ID:             "668d797c-76fa-4934-9b05-ad288df2d136",
						CreatedTime:    time.Date(2020, 3, 17, 19, 10, 4, 968_000_000, time.UTC),
						LastEditedTime: time.Date(2020, 3, 17, 21, 49, 37, 913_000_000, time.UTC),
						Title: []RichText{
							&RichTextText{
								BaseRichText: BaseRichText{
									PlainText: "Grocery List",
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
									Content: "Grocery List",
									Link:    nil,
								},
							},
						},
						Properties: map[string]Property{
							"Name": &TitleProperty{
								baseProperty: baseProperty{
									ID:   "title",
									Type: "title",
								},
								Title: map[string]interface{}{},
							},
							// FIXME: This is different from the example of https://developers.notion.com/reference/get-database
							// but in the API reference there's no `text` type, thus changing it to `rich_text`, but need
							// to confirm what is the expected type.
							"Description": &RichTextProperty{
								baseProperty: baseProperty{
									ID:   "J@cS",
									Type: PropertyTypeRichText,
								},
								RichText: map[string]interface{}{},
							},
							"In stock": &CheckboxProperty{
								baseProperty: baseProperty{
									ID:   "{xY`",
									Type: PropertyTypeCheckbox,
								},
								Checkbox: map[string]interface{}{},
							},
							"Food group": &SelectProperty{
								baseProperty: baseProperty{
									ID:   "TJmr",
									Type: PropertyTypeSelect,
								},
								Select: SelectPropertyOption{
									Options: []SelectOption{
										{
											ID:    "96eb622f-4b88-4283-919d-ece2fbed3841",
											Name:  "🥦Vegetable",
											Color: ColorGreen,
										},
										{
											ID:    "bb443819-81dc-46fb-882d-ebee6e22c432",
											Name:  "🍎Fruit",
											Color: ColorRed,
										},
										{
											ID:    "7da9d1b9-8685-472e-9da3-3af57bdb221e",
											Name:  "💪Protein",
											Color: ColorYellow,
										},
									},
								},
							},
							"Price": &NumberProperty{
								baseProperty: baseProperty{
									ID:   "cU^N",
									Type: PropertyTypeNumber,
								},
								Number: NumberPropertyOption{
									Format: NumberFormatDollar,
								},
							},
							"Cost of next trip": &FormulaProperty{
								baseProperty: baseProperty{
									ID:   "p:sC",
									Type: PropertyTypeFormula,
								},
								Formula: Formula{
									// FIXME: The example response in https://developers.notion.com/reference/get-database
									// has the key `value` but in the API reference https://developers.notion.com/reference/database#formula-configuration
									// the property name should be `expression`, need to check if this is expected.
									Expression: `if(prop("In stock"), 0, prop("Price"))`,
								},
							},
							"Last ordered": &DateProperty{
								baseProperty: baseProperty{
									ID:   "]\\R[",
									Type: PropertyTypeDate,
								},
								Date: map[string]interface{}{},
							},
							"Meals": &RelationProperty{
								baseProperty: baseProperty{
									// FIXME: The example response in https://developers.notion.com/reference/get-database
									// does not contain the ID, need to check if this is expected.
									ID:   "",
									Type: PropertyTypeRelation,
								},
								Relation: Relation{
									// FIXME: The key in the example is `database` but should be `database_id` instead,
									// and need to check if which one is correct.
									DatabaseID: "668d797c-76fa-4934-9b05-ad288df2d136",
								},
							},
							"Number of meals": &RollupProperty{
								baseProperty: baseProperty{
									ID:   "Z\\Eh",
									Type: PropertyTypeRollup,
								},
								Rollup: RollupPropertyOption{
									RelationPropertyName: "Meals",
									RelationPropertyID:   "mxp^",
									RollupPropertyName:   "Name",
									RollupPropertyID:     "title",
									// FIXME: In the example the returned function is `count` but the possible values
									// do not include `count`, thus need to confirm this.
									Function: "count",
								},
							},
							"Store availability": &MultiSelectProperty{
								baseProperty: baseProperty{
									// FIXME: The example response in https://developers.notion.com/reference/get-database
									// does not contain the ID, need to check if this is expected.
									ID:   "",
									Type: PropertyTypeMultiSelect,
								},
								MultiSelect: MultiSelectPropertyOption{
									Options: []MultiSelectOption{
										{
											ID:    "d209b920-212c-4040-9d4a-bdf349dd8b2a",
											Name:  "Duc Loi Market",
											Color: ColorBlue,
										},
										{
											ID:    "70104074-0f91-467b-9787-00d59e6e1e41",
											Name:  "Rainbow Grocery",
											Color: ColorGray,
										},
										{
											ID:    "e6fd4f04-894d-4fa7-8d8b-e92d08ebb604",
											Name:  "Nijiya Market",
											Color: ColorPurple,
										},
										{
											ID:    "6c3867c5-d542-4f84-b6e9-a420c43094e7",
											Name:  "Gus's Community Market",
											Color: ColorYellow,
										},
									},
								},
							},
							"+1": &PeopleProperty{
								baseProperty: baseProperty{
									ID:   "aGut",
									Type: PropertyTypePeople,
								},
								People: map[string]interface{}{},
							},
							"Photo": &FileProperty{
								baseProperty: baseProperty{
									ID:   "aTIT",
									Type: PropertyTypeFile,
								},
								File: map[string]interface{}{},
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

			d := &databasesClient{
				restClient: tt.fields.restClient.BaseURL(mockHTTPServer.URL),
			}
			got, err := d.Retrieve(tt.args.ctx, tt.args.params)
			if tt.wants.err != nil {
				assert.ErrorIs(t, err, tt.wants.err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.wants.response, got)
		})
	}
}

func newFloat64(f float64) *float64 {
	return &f
}
