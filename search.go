package notion

import (
	"context"
	"encoding/json"

	"github.com/mkfsn/notion-go/rest"
	"github.com/mkfsn/notion-go/typed"
)

type SearchFilter struct {
	// The value of the property to filter the results by. Possible values for object type include `page` or `database`.
	// Limitation: Currently the only filter allowed is object which will filter by type of `object`
	// (either `page` or `database`)
	Value typed.SearchFilterValue `json:"value"`
	// The name of the property to filter by. Currently the only property you can filter by is the object type.
	// Possible values include `object`. Limitation: Currently the only filter allowed is `object` which will
	// filter by type of object (either `page` or `database`)
	Property typed.SearchFilterProperty `json:"property"`
}

type SearchSort struct {
	// The direction to sort.
	Direction typed.SearchSortDirection `json:"direction"`
	// The name of the timestamp to sort against. Possible values include `last_edited_time`.
	Timestamp typed.SearchSortTimestamp `json:"timestamp"`
}

type SearchParameters struct {
	PaginationParameters
	Query  string       `json:"query" url:"-"`
	Sort   SearchSort   `json:"sort" url:"-"`
	Filter SearchFilter `json:"filter" url:"-"`
}

type SearchableObject interface {
	isSearchable()
}

type SearchResponse struct {
	PaginatedList
	Results []SearchableObject `json:"results"`
}

func (s *SearchResponse) UnmarshalJSON(data []byte) error {
	type Alias SearchResponse

	alias := struct {
		*Alias
		Results []searchableObjectDecoder `json:"results"`
	}{
		Alias: (*Alias)(s),
	}

	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}

	s.Results = make([]SearchableObject, 0, len(alias.Results))

	for _, decoder := range alias.Results {
		s.Results = append(s.Results, decoder.SearchableObject)
	}

	return nil
}

type SearchInterface interface {
	Search(ctx context.Context, params SearchParameters) (*SearchResponse, error)
}

type searchClient struct {
	restClient rest.Interface
}

func newSearchClient(restClient rest.Interface) *searchClient {
	return &searchClient{
		restClient: restClient,
	}
}

func (s *searchClient) Search(ctx context.Context, params SearchParameters) (*SearchResponse, error) {
	var result SearchResponse
	var failure HTTPError

	err := s.restClient.New().
		Post().
		Endpoint(APISearchEndpoint).
		QueryStruct(params).
		BodyJSON(params).
		Receive(ctx, &result, &failure)

	return &result, err
}

type searchableObjectDecoder struct {
	SearchableObject
}

func (s *searchableObjectDecoder) UnmarshalJSON(data []byte) error {
	var decoder struct {
		Object typed.ObjectType `json:"object"`
	}

	if err := json.Unmarshal(data, &decoder); err != nil {
		return err
	}

	switch decoder.Object {
	case typed.ObjectTypePage:
		s.SearchableObject = &Page{}

	case typed.ObjectTypeDatabase:
		s.SearchableObject = &Database{}

	case typed.ObjectTypeBlock:
		return ErrUnknown
	}

	return json.Unmarshal(data, s.SearchableObject)
}
