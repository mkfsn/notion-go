package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/google/go-querystring/query"
)

type restClient struct {
	baseURL    string
	header     http.Header
	httpClient *http.Client

	method      string
	endpoint    string
	queryStruct interface{}
	bodyJSON    interface{}
}

func New() Interface {
	return &restClient{
		header:     make(http.Header),
		httpClient: http.DefaultClient,
	}
}

func (r *restClient) New() Interface {
	newRestClient := &restClient{
		baseURL:    r.baseURL,
		header:     r.header.Clone(),
		httpClient: r.httpClient, // TODO: deep copy
	}
	return newRestClient
}

func (r *restClient) BearerToken(token string) Interface {
	r.header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	return r
}

func (r *restClient) BaseURL(baseURL string) Interface {
	r.baseURL = baseURL
	return r
}

func (r *restClient) Client(httpClient *http.Client) Interface {
	r.httpClient = httpClient
	return r
}

func (r *restClient) UserAgent(userAgent string) Interface {
	r.header.Set("User-Agent", userAgent)
	return r
}

func (r *restClient) Header(key, value string) Interface {
	r.header.Set(key, value)
	return r
}

func (r *restClient) Get() Interface {
	r.method = http.MethodGet
	return r
}

func (r *restClient) Post() Interface {
	r.method = http.MethodPost
	return r
}

func (r *restClient) Patch() Interface {
	r.method = http.MethodPatch
	return r
}

func (r *restClient) Endpoint(endpoint string) Interface {
	r.endpoint = endpoint
	return r
}

func (r *restClient) QueryStruct(queryStruct interface{}) Interface {
	r.queryStruct = queryStruct
	return r
}

func (r *restClient) BodyJSON(bodyJSON interface{}) Interface {
	r.bodyJSON = bodyJSON

	if r.bodyJSON != nil {
		r.header.Add("Content-Type", "application/json")
	}

	return r
}

func (r *restClient) Request(ctx context.Context) (*http.Request, error) {
	v, err := query.Values(r.queryStruct)
	if err != nil {
		return nil, err
	}

	b, err := json.Marshal(r.bodyJSON)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, r.method, r.baseURL+r.endpoint, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = v.Encode()

	req.Header = r.header

	return req, nil
}

func (r *restClient) Receive(ctx context.Context, success, failure interface{}) error {
	req, err := r.Request(ctx)
	if err != nil {
		return err
	}

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return r.decodeResponseData(resp.StatusCode, b, success, failure)
}

func (r *restClient) decodeResponseData(statusCode int, data []byte, success, failure interface{}) error {
	if statusCode == http.StatusOK {
		if success != nil {
			return json.Unmarshal(data, success)
		}

		return nil
	}

	if failure == nil {
		return nil
	}

	if err := json.Unmarshal(data, failure); err != nil {
		return err
	}

	return failure.(error)
}
