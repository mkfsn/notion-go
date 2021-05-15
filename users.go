package notion

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

type UserType = string

const (
	UserTypePerson UserType = "person"
	UserTypeBot    UserType = "bot"
)

type User interface {
	isUser()
}

type baseUser struct {
	Object    string   `json:"object"`
	ID        string   `json:"id"`
	Type      UserType `json:"type"`
	Name      string   `json:"name"`
	AvatarURL string   `json:"avatar_url"`
}

func (b baseUser) isUser() {}

type PersonUser struct {
	baseUser
	Person *struct {
		Email string `json:"email"`
	} `json:"person"`
}

type BotUser struct {
	baseUser
	Bot interface{} `json:"bot"`
}

type UsersRetrieveParameters struct {
	UserID string `json:"user_id"`
}

type UsersRetrieveResponse struct {
	User
}

func (u *UsersRetrieveResponse) UnmarshalJSON(data []byte) error {
	var base baseUser

	if err := json.Unmarshal(data, &base); err != nil {
		return err
	}

	switch base.Type {
	case UserTypePerson:
		var user PersonUser

		if err := json.Unmarshal(data, &user); err != nil {
			return err
		}

		u.User = user

	case UserTypeBot:
		var user BotUser

		if err := json.Unmarshal(data, &user); err != nil {
			return err
		}

		u.User = user
	}

	return nil
}

type UsersListParameters struct {
	StartCursor string `json:"start_cursor"`
	PageSize    int32  `json:"page_size"`
}

type UsersListResponse struct {
	Results    []User `json:"results"`
	NextCursor string `json:"next_cursor"`
	HasMore    bool   `json:"has_more"`
}

func (u *UsersListResponse) UnmarshalJSON(data []byte) error {
	type Alias UsersListResponse

	alias := struct {
		*Alias
		Results []json.RawMessage `json:"results"`
	}{
		Alias: (*Alias)(u),
	}

	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}

	u.Results = make([]User, 0, len(alias.Results))

	for _, result := range alias.Results {
		var base baseUser

		if err := json.Unmarshal(result, &base); err != nil {
			return err
		}

		switch base.Type {
		case UserTypePerson:
			var user PersonUser

			if err := json.Unmarshal(result, &user); err != nil {
				return err
			}

			u.Results = append(u.Results, user)

		case UserTypeBot:
			var user BotUser

			if err := json.Unmarshal(result, &user); err != nil {
				return err
			}

			u.Results = append(u.Results, user)
		}
	}

	return nil
}

type UsersInterface interface {
	Retrieve(ctx context.Context, params UsersRetrieveParameters) (*UsersRetrieveResponse, error)
	List(ctx context.Context, params UsersListParameters) (*UsersListResponse, error)
}

type usersClient struct {
	client client
}

func newUsersClient(client client) *usersClient {
	return &usersClient{
		client: client,
	}
}

func (u *usersClient) Retrieve(ctx context.Context, params UsersRetrieveParameters) (*UsersRetrieveResponse, error) {
	endpoint := strings.Replace(APIUsersRetrieveEndpoint, "{user_id}", params.UserID, 1)

	b, err := u.client.Request(ctx, http.MethodGet, endpoint, params)
	if err != nil {
		return nil, err
	}

	var response UsersRetrieveResponse

	if err := json.Unmarshal(b, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (u *usersClient) List(ctx context.Context, params UsersListParameters) (*UsersListResponse, error) {
	b, err := u.client.Request(ctx, http.MethodGet, APIUsersListEndpoint, params)
	if err != nil {
		return nil, err
	}

	var response UsersListResponse

	if err := json.Unmarshal(b, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
