package rest

import (
	"context"
	"net/http"
)

type Interface interface {
	New() Interface
	BearerToken(token string) Interface
	BaseURL(baseURL string) Interface
	Client(httpClient *http.Client) Interface
	UserAgent(userAgent string) Interface
	Header(key, value string) Interface
	Get() Interface
	Post() Interface
	Patch() Interface
	Endpoint(endpoint string) Interface
	QueryStruct(queryStruct interface{}) Interface
	BodyJSON(bodyJSON interface{}) Interface
	Request(ctx context.Context) (*http.Request, error)
	Receive(ctx context.Context, success, failure interface{}) error
}
