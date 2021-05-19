package notion

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/google/go-querystring/query"
)

var (
	defaultOptions = clientOptions{
		authToken:     "",
		baseURL:       APIBaseURL,
		notionVersion: DefaultNotionVersion,
	}
)

type clientOptions struct {
	authToken     string
	baseURL       string
	notionVersion string
}

type ClientOption func(o *clientOptions)

func WithAuthToken(authToken string) ClientOption {
	return func(o *clientOptions) {
		o.authToken = authToken
	}
}

func WithBaseURL(baseURL string) ClientOption {
	return func(o *clientOptions) {
		o.baseURL = baseURL
	}
}

func WithNotionVersion(notionVersion string) ClientOption {
	return func(o *clientOptions) {
		o.notionVersion = notionVersion
	}
}

type client interface {
	Request(ctx context.Context, method string, endpoint string, body interface{}) ([]byte, error)
}

type httpClient struct {
	http.Client
	clientOptions
}

func newHTTPClient(clientOptions clientOptions) *httpClient {
	return &httpClient{
		clientOptions: clientOptions,
	}
}

func (c *httpClient) Request(ctx context.Context, method string, endpoint string, body interface{}) ([]byte, error) {
	v, err := query.Values(body)
	if err != nil {
		return nil, err
	}

	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+endpoint, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = v.Encode()

	if body != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	req.Header.Add("Notion-Version", c.notionVersion)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.authToken))

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, newHTTPError(resp.StatusCode, b)
	}

	return b, nil
}
