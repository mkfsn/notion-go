package notion

import (
	"context"
	"encoding/json"
	"net/http"

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
		Results []json.RawMessage `json:"results"`
	}{
		Alias: (*Alias)(s),
	}

	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}

	s.Results = make([]SearchableObject, 0, len(alias.Results))

	for _, result := range alias.Results {
		var base struct {
			Object typed.ObjectType `json:"object"`
		}

		if err := json.Unmarshal(result, &base); err != nil {
			return err
		}

		switch base.Object {
		case typed.ObjectTypePage:
			var object Page

			if err := json.Unmarshal(result, &object); err != nil {
				return err
			}

			s.Results = append(s.Results, object)

		case typed.ObjectTypeDatabase:
			var object Database

			if err := json.Unmarshal(result, &object); err != nil {
				return err
			}

			s.Results = append(s.Results, object)

		case typed.ObjectTypeBlock:
			continue
		}
	}

	return nil
}

type SearchInterface interface {
	Search(ctx context.Context, params SearchParameters) (*SearchResponse, error)
}

type searchClient struct {
	client client
}

func newSearchClient(client client) *searchClient {
	return &searchClient{
		client: client,
	}
}

func (s *searchClient) Search(ctx context.Context, params SearchParameters) (*SearchResponse, error) {
	b, err := s.client.Request(ctx, http.MethodPost, APISearchEndpoint, params)
	if err != nil {
		return nil, err
	}

	var response SearchResponse

	if err := json.Unmarshal(b, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
