package notion

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/mkfsn/notion-go/rest"
	"github.com/mkfsn/notion-go/typed"
)

type User interface {
	isUser()
}

func newUser(data []byte) (User, error) {
	var base baseUser

	if err := json.Unmarshal(data, &base); err != nil {
		return nil, err
	}

	switch base.Type {
	case typed.UserTypePerson:
		var user PersonUser

		if err := json.Unmarshal(data, &user); err != nil {
			return nil, err
		}

		return user, nil

	case typed.UserTypeBot:
		var user BotUser

		if err := json.Unmarshal(data, &user); err != nil {
			return nil, err
		}

		return user, nil
	}

	return nil, ErrUnknown
}

type baseUser struct {
	Object    string         `json:"object"`
	ID        string         `json:"id"`
	Type      typed.UserType `json:"type"`
	Name      string         `json:"name"`
	AvatarURL string         `json:"avatar_url"`
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
	UserID string `json:"-"`
}

type UsersRetrieveResponse struct {
	User
}

func (u *UsersRetrieveResponse) UnmarshalJSON(data []byte) (err error) {
	u.User, err = newUser(data)
	return
}

type UsersListParameters struct {
	PaginationParameters
}

type UsersListResponse struct {
	PaginatedList
	Results []User `json:"results"`
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
		user, err := newUser(result)
		if err != nil {
			return err
		}

		u.Results = append(u.Results, user)
	}

	return nil
}

type UsersInterface interface {
	Retrieve(ctx context.Context, params UsersRetrieveParameters) (*UsersRetrieveResponse, error)
	List(ctx context.Context, params UsersListParameters) (*UsersListResponse, error)
}

type usersClient struct {
	restClient rest.Interface
}

func newUsersClient(restClient rest.Interface) *usersClient {
	return &usersClient{
		restClient: restClient,
	}
}

func (u *usersClient) Retrieve(ctx context.Context, params UsersRetrieveParameters) (*UsersRetrieveResponse, error) {
	var result UsersRetrieveResponse
	var failure HTTPError

	err := u.restClient.New().Get().
		Endpoint(strings.Replace(APIUsersRetrieveEndpoint, "{user_id}", params.UserID, 1)).
		QueryStruct(params).
		BodyJSON(nil).
		Receive(ctx, &result, &failure)

	return &result, err
}

func (u *usersClient) List(ctx context.Context, params UsersListParameters) (*UsersListResponse, error) {
	var result UsersListResponse
	var failure HTTPError

	err := u.restClient.New().Get().
		Endpoint(APIUsersListEndpoint).
		QueryStruct(params).
		BodyJSON(params).
		Receive(ctx, &result, &failure)

	return &result, err
}
