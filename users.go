package notion

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/mkfsn/notion-go/rest"
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

type Person struct {
	Email string `json:"email"`
}

type PersonUser struct {
	baseUser
	Person Person `json:"person"`
}

type Bot struct{}

type BotUser struct {
	baseUser
	Bot Bot `json:"bot"`
}

type UsersRetrieveParameters struct {
	UserID string `json:"-" url:"-"`
}

type UsersRetrieveResponse struct {
	User
}

func (u *UsersRetrieveResponse) UnmarshalJSON(data []byte) (err error) {
	var decoder userDecoder

	if err := json.Unmarshal(data, &decoder); err != nil {
		return err
	}

	u.User = decoder.User

	return nil
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
		Results []userDecoder `json:"results"`
	}{
		Alias: (*Alias)(u),
	}

	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}

	u.Results = make([]User, 0, len(alias.Results))

	for _, decoder := range alias.Results {
		u.Results = append(u.Results, decoder.User)
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

type userDecoder struct {
	User
}

func (u *userDecoder) UnmarshalJSON(data []byte) error {
	var decoder struct {
		Type UserType `json:"type"`
	}

	if err := json.Unmarshal(data, &decoder); err != nil {
		return err
	}

	switch decoder.Type {
	case UserTypePerson:
		u.User = &PersonUser{}

	case UserTypeBot:
		u.User = &BotUser{}
	}

	return json.Unmarshal(data, u.User)
}
