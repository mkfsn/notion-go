package notion

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mkfsn/notion-go/rest"
	"github.com/stretchr/testify/assert"
)

func Test_usersClient_Retrieve(t *testing.T) {
	type fields struct {
		restClient      rest.Interface
		mockHTTPHandler http.Handler
		authToken       string
	}

	type args struct {
		ctx    context.Context
		params UsersRetrieveParameters
	}

	type wants struct {
		response *UsersRetrieveResponse
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
			name: "Bot User",
			fields: fields{
				restClient: rest.New(),
				authToken:  "033dcdcf-8252-49f4-826c-e795fcab0ad2",
				mockHTTPHandler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
					assert.Equal(t, DefaultNotionVersion, request.Header.Get("Notion-Version"))
					assert.Equal(t, DefaultUserAgent, request.Header.Get("User-Agent"))
					assert.Equal(t, "Bearer 033dcdcf-8252-49f4-826c-e795fcab0ad2", request.Header.Get("Authorization"))

					assert.Equal(t, http.MethodGet, request.Method)
					assert.Equal(t, "/v1/users/9a3b5ae0-c6e6-482d-b0e1-ed315ee6dc57", request.RequestURI)

					writer.WriteHeader(http.StatusOK)

					_, err := writer.Write([]byte(`{
    					  "object": "user",
    					  "id": "9a3b5ae0-c6e6-482d-b0e1-ed315ee6dc57",
    					  "type": "bot",
    					  "bot": {},
    					  "name": "Doug Engelbot",
    					  "avatar_url": "https://secure.notion-static.com/6720d746-3402-4171-8ebb-28d15144923c.jpg"
    					}`,
					))
					assert.NoError(t, err)
				}),
			},
			args: args{
				ctx:    context.Background(),
				params: UsersRetrieveParameters{UserID: "9a3b5ae0-c6e6-482d-b0e1-ed315ee6dc57"},
			},
			wants: wants{
				response: &UsersRetrieveResponse{
					User: &BotUser{
						baseUser: baseUser{
							Object:    ObjectTypeUser,
							ID:        "9a3b5ae0-c6e6-482d-b0e1-ed315ee6dc57",
							Type:      UserTypeBot,
							Name:      "Doug Engelbot",
							AvatarURL: "https://secure.notion-static.com/6720d746-3402-4171-8ebb-28d15144923c.jpg",
						},
						Bot: Bot{},
					},
				},
			},
		},

		{
			name: "Person User",
			fields: fields{
				restClient: rest.New(),
				authToken:  "033dcdcf-8252-49f4-826c-e795fcab0ad2",
				mockHTTPHandler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
					assert.Equal(t, DefaultNotionVersion, request.Header.Get("Notion-Version"))
					assert.Equal(t, DefaultUserAgent, request.Header.Get("User-Agent"))
					assert.Equal(t, "Bearer 033dcdcf-8252-49f4-826c-e795fcab0ad2", request.Header.Get("Authorization"))

					assert.Equal(t, http.MethodGet, request.Method)
					assert.Equal(t, "/v1/users/d40e767c-d7af-4b18-a86d-55c61f1e39a4", request.RequestURI)

					writer.WriteHeader(http.StatusOK)

					_, err := writer.Write([]byte(`{
    					  "object": "user",
    					  "id": "d40e767c-d7af-4b18-a86d-55c61f1e39a4",
    					  "type": "person",
    					  "person": {
    					    "email": "avo@example.org"
    					  },
    					  "name": "Avocado Lovelace",
    					  "avatar_url": "https://secure.notion-static.com/e6a352a8-8381-44d0-a1dc-9ed80e62b53d.jpg"
    					}`,
					))
					assert.NoError(t, err)
				}),
			},
			args: args{
				ctx:    context.Background(),
				params: UsersRetrieveParameters{UserID: "d40e767c-d7af-4b18-a86d-55c61f1e39a4"},
			},
			wants: wants{
				response: &UsersRetrieveResponse{
					User: &PersonUser{
						baseUser: baseUser{
							Object:    ObjectTypeUser,
							ID:        "d40e767c-d7af-4b18-a86d-55c61f1e39a4",
							Type:      UserTypePerson,
							Name:      "Avocado Lovelace",
							AvatarURL: "https://secure.notion-static.com/e6a352a8-8381-44d0-a1dc-9ed80e62b53d.jpg",
						},
						Person: Person{Email: "avo@example.org"},
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

			got, err := sut.Users().Retrieve(tt.args.ctx, tt.args.params)
			if tt.wants.err != nil {
				assert.ErrorIs(t, err, tt.wants.err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.wants.response, got)
		})
	}
}

func Test_usersClient_List(t *testing.T) {
	type fields struct {
		restClient      rest.Interface
		mockHTTPHandler http.Handler
		authToken       string
	}

	type args struct {
		ctx    context.Context
		params UsersListParameters
	}

	type wants struct {
		response *UsersListResponse
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
			name: "List two users in one page",
			fields: fields{
				restClient: rest.New(),
				authToken:  "2a966a7a-6e97-4b2c-abb2-c0eba4dbcb5f",
				mockHTTPHandler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
					assert.Equal(t, DefaultNotionVersion, request.Header.Get("Notion-Version"))
					assert.Equal(t, DefaultUserAgent, request.Header.Get("User-Agent"))
					assert.Equal(t, "Bearer 2a966a7a-6e97-4b2c-abb2-c0eba4dbcb5f", request.Header.Get("Authorization"))

					assert.Equal(t, http.MethodGet, request.Method)
					assert.Equal(t, "/v1/users?page_size=2", request.RequestURI)

					writer.WriteHeader(http.StatusOK)

					_, err := writer.Write([]byte(`{
					      "results": [
					        {
					          "object": "user",
					          "id": "d40e767c-d7af-4b18-a86d-55c61f1e39a4",
					          "type": "person",
					          "person": {
					            "email": "avo@example.org"
					          },
					          "name": "Avocado Lovelace",
					          "avatar_url": "https://secure.notion-static.com/e6a352a8-8381-44d0-a1dc-9ed80e62b53d.jpg"
					        },
					        {
					          "object": "user",
					          "id": "9a3b5ae0-c6e6-482d-b0e1-ed315ee6dc57",
					          "type": "bot",
					          "bot": {},
					          "name": "Doug Engelbot",
					          "avatar_url": "https://secure.notion-static.com/6720d746-3402-4171-8ebb-28d15144923c.jpg"
					        }
					      ],
					      "next_cursor": "fe2cc560-036c-44cd-90e8-294d5a74cebc",
					      "has_more": true
						}`,
					))
					assert.NoError(t, err)
				}),
			},
			args: args{
				ctx: context.Background(),
				params: UsersListParameters{
					PaginationParameters: PaginationParameters{
						StartCursor: "",
						PageSize:    2,
					},
				},
			},
			wants: wants{
				response: &UsersListResponse{
					PaginatedList: PaginatedList{
						NextCursor: "fe2cc560-036c-44cd-90e8-294d5a74cebc",
						HasMore:    true,
					},
					Results: []User{
						&PersonUser{
							baseUser: baseUser{
								Object:    ObjectTypeUser,
								ID:        "d40e767c-d7af-4b18-a86d-55c61f1e39a4",
								Type:      UserTypePerson,
								Name:      "Avocado Lovelace",
								AvatarURL: "https://secure.notion-static.com/e6a352a8-8381-44d0-a1dc-9ed80e62b53d.jpg",
							},
							Person: Person{Email: "avo@example.org"},
						},
						&BotUser{
							baseUser: baseUser{
								Object:    ObjectTypeUser,
								ID:        "9a3b5ae0-c6e6-482d-b0e1-ed315ee6dc57",
								Type:      UserTypeBot,
								Name:      "Doug Engelbot",
								AvatarURL: "https://secure.notion-static.com/6720d746-3402-4171-8ebb-28d15144923c.jpg",
							},
							Bot: Bot{},
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

			sut := New(
				tt.fields.authToken,
				WithBaseURL(mockHTTPServer.URL),
			)

			got, err := sut.Users().List(tt.args.ctx, tt.args.params)
			if tt.wants.err != nil {
				assert.ErrorIs(t, err, tt.wants.err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.wants.response, got)
		})
	}
}
