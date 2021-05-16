package notion

import (
	"github.com/mkfsn/notion-go/typed"
)

type PaginationParameters struct {
	// If supplied, this endpoint will return a page of results starting after the cursor provided.
	// If not supplied, this endpoint will return the first page of results.
	StartCursor string `json:"-" url:"start_cursor,omitempty"`
	// The number of items from the full list desired in the response. Maximum: 100
	PageSize int32 `json:"-" url:"page_size,omitempty"`
}

type PaginatedList struct {
	Object     typed.ObjectType `json:"object"`
	HasMore    bool             `json:"has_more"`
	NextCursor string           `json:"next_cursor"`
}
