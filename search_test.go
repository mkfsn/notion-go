package notion

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mkfsn/notion-go/rest"
	"github.com/stretchr/testify/assert"
)

func Test_searchClient_Search(t *testing.T) {
	type fields struct {
		restClient      rest.Interface
		mockHTTPHandler http.Handler
	}

	type args struct {
		ctx    context.Context
		params SearchParameters
	}

	type wants struct {
		response *SearchResponse
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
			name: "Search two objects in one page",
			fields: fields{
				restClient: rest.New(),
				mockHTTPHandler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
					assert.Equal(t, http.MethodPost, request.Method)
					assert.Equal(t, "/v1/search?page_size=2", request.RequestURI)

					writer.WriteHeader(http.StatusOK)

					_, err := writer.Write([]byte(`{
					    "has_more": false,
					    "next_cursor": null,
					    "object": "list",
					    "results": [
							{
					            "created_time": "2021-04-22T22:23:26.080Z",
					            "id": "e6c6f8ff-c70e-4970-91ba-98f03e0d7fc6",
					            "last_edited_time": "2021-04-23T04:21:00.000Z",
					            "object": "database",
					            "properties": {
					                "Name": {
					                    "id": "title",
					                    "title": {},
					                    "type": "title"
					                },
					                "Task Type": {
					                    "id": "vd@l",
					                    "multi_select": {
					                        "options": []
					                    },
					                    "type": "multi_select"
					                }
					            },
					            "title": [
					                {
					                    "annotations": {
					                        "bold": false,
					                        "code": false,
					                        "color": "default",
					                        "italic": false,
					                        "strikethrough": false,
					                        "underline": false
					                    },
					                    "href": null,
					                    "plain_text": "Tasks",
					                    "text": {
					                        "content": "Tasks",
					                        "link": null
					                    },
					                    "type": "text"
					                }
					            ]
					        },
							{
					            "archived": false,
					            "created_time": "2021-04-23T04:21:00.000Z",
					            "id": "4f555b50-3a9b-49cb-924c-3746f4ca5522",
					            "last_edited_time": "2021-04-23T04:21:00.000Z",
					            "object": "page",
					            "parent": {
					                "database_id": "e6c6f8ff-c70e-4970-91ba-98f03e0d7fc6",
					                "type": "database_id"
					            },
					            "properties": {
					                "Name": {
					                    "id": "title",
					                    "title": [
					                        {
					                            "annotations": {
					                                "bold": false,
					                                "code": false,
					                                "color": "default",
					                                "italic": false,
					                                "strikethrough": false,
					                                "underline": false
					                            },
					                            "href": null,
					                            "plain_text": "Task 1",
					                            "text": {
					                                "content": "Task1 1",
					                                "link": null
					                            },
					                            "type": "text"
					                        }
					                    ],
					                    "type": "title"
					                }
					            }
					        }
					    ]
					}`,
					))
					assert.NoError(t, err)
				}),
			},
			args: args{
				ctx: context.Background(),
				params: SearchParameters{
					PaginationParameters: PaginationParameters{
						StartCursor: "",
						PageSize:    2,
					},
					Query: "External tasks",
					Sort: SearchSort{
						Direction: SearchSortDirectionAscending,
						Timestamp: SearchSortTimestampLastEditedTime,
					},
					Filter: SearchFilter{},
				},
			},
			wants: wants{
				response: &SearchResponse{
					PaginatedList: PaginatedList{
						Object:     ObjectTypeList,
						HasMore:    false,
						NextCursor: "",
					},
					Results: []SearchableObject{
						&Database{
							Object:         ObjectTypeDatabase,
							ID:             "e6c6f8ff-c70e-4970-91ba-98f03e0d7fc6",
							CreatedTime:    time.Date(2021, 4, 22, 22, 23, 26, 80_000_000, time.UTC),
							LastEditedTime: time.Date(2021, 4, 23, 4, 21, 0, 0, time.UTC),
							Title: []RichText{
								&RichTextText{
									BaseRichText: BaseRichText{
										PlainText: "Tasks",
										Href:      "",
										Type:      RichTextTypeText,
										Annotations: &Annotations{
											Bold:          false,
											Italic:        false,
											Strikethrough: false,
											Underline:     false,
											Code:          false,
											Color:         DefaultColor,
										},
									},
									Text: TextObject{
										Content: "Tasks",
										Link:    nil,
									},
								},
							},
							Properties: map[string]Property{
								"Name": &TitleProperty{
									baseProperty: baseProperty{
										ID:   "title",
										Type: PropertyTypeTitle,
									},
									Title: map[string]interface{}{},
								},
								"Task Type": &MultiSelectProperty{
									baseProperty: baseProperty{
										ID:   "vd@l",
										Type: PropertyTypeMultiSelect,
									},
									MultiSelect: MultiSelectPropertyOption{
										Options: []MultiSelectOption{},
									},
								},
							},
						},
						&Page{
							Object: ObjectTypePage,
							ID:     "4f555b50-3a9b-49cb-924c-3746f4ca5522",
							Parent: &DatabaseParent{
								baseParent: baseParent{
									Type: ParentTypeDatabase,
								},
								DatabaseID: "e6c6f8ff-c70e-4970-91ba-98f03e0d7fc6",
							},
							Properties: map[string]PropertyValue{
								"Name": &TitlePropertyValue{
									basePropertyValue: basePropertyValue{
										ID:   "title",
										Type: PropertyValueTypeTitle,
									},
									Title: []RichText{

										&RichTextText{
											BaseRichText: BaseRichText{
												PlainText: "Task 1",
												Href:      "",
												Type:      RichTextTypeText,
												Annotations: &Annotations{
													Bold:          false,
													Italic:        false,
													Strikethrough: false,
													Underline:     false,
													Code:          false,
													Color:         DefaultColor,
												},
											},
											Text: TextObject{
												Content: "Task1 1",
												Link:    nil,
											},
										},
									},
								},
							},
							CreatedTime:    time.Date(2021, 4, 23, 4, 21, 0, 0, time.UTC),
							LastEditedTime: time.Date(2021, 4, 23, 4, 21, 0, 0, time.UTC),
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

			sut := &searchClient{
				restClient: tt.fields.restClient.BaseURL(mockHTTPServer.URL),
			}

			got, err := sut.Search(tt.args.ctx, tt.args.params)
			if tt.wants.err != nil {
				assert.ErrorIs(t, err, tt.wants.err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.wants.response, got)
		})
	}
}
