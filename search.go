package notion

import (
	"context"
	"encoding/json"
	"net/http"
)

type SearchFilterValue string

const (
	SearchFilterValuePage     SearchFilterValue = "page"
	SearchFilterValueDatabase SearchFilterValue = "database"
)

type SearchFilterProperty string

const (
	SearchFilterPropertyObject SearchFilterProperty = "object"
)

type SearchFilter struct {
	// The value of the property to filter the results by. Possible values for object type include `page` or `database`.
	// Limitation: Currently the only filter allowed is object which will filter by type of `object`
	// (either `page` or `database`)
	Value SearchFilterValue `json:"value"`
	// The name of the property to filter by. Currently the only property you can filter by is the object type.
	// Possible values include `object`. Limitation: Currently the only filter allowed is `object` which will
	// filter by type of object (either `page` or `database`)
	Property SearchFilterProperty `json:"property"`
}

type SearchSortDirection string

const (
	SearchSortDirectionAscending  SearchSortDirection = "ascending"
	SearchSortDirectionDescending SearchSortDirection = " descending"
)

type SearchSortTimestamp string

const (
	SearchSortTimestampLastEditedTime SearchSortTimestamp = "last_edited_time"
)

type SearchSort struct {
	// The direction to sort.
	Direction SearchSortDirection `json:"direction"`
	// The name of the timestamp to sort against. Possible values include `last_edited_time`.
	Timestamp SearchSortTimestamp `json:"timestamp"`
}

type SearchParameters struct {
	PaginationParameters
	Query  string       `json:"query"`
	Sort   SearchSort   `json:"sort"`
	Filter SearchFilter `json:"filter"`
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
			Object ObjectType `json:"object"`
		}

		if err := json.Unmarshal(result, &base); err != nil {
			return err
		}

		switch base.Object {
		case ObjectTypePage:
			var object Page

			if err := json.Unmarshal(result, &object); err != nil {
				return err
			}

			s.Results = append(s.Results, object)

		case ObjectTypeDatabase:
			var object Database

			if err := json.Unmarshal(result, &object); err != nil {
				return err
			}

			s.Results = append(s.Results, object)
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
