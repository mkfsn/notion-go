package notion

import (
	"context"
)

const (
	APIBaseURL                      = "https://api.notion.com"
	APIUsersListEndpoint            = "/v1/users"
	APIUsersRetrieveEndpoint        = "/v1/users/{user_id}"
	APIBlocksListChildrenEndpoint   = "/v1/blocks/{block_id}/children"
	APIBlocksAppendChildrenEndpoint = "/v1/blocks/{block_id}/children"
	APIPagesCreateEndpoint          = "/v1/pages"
	APIPagesRetrieveEndpoint        = "/v1/pages/{page_id}"
	APIPagesUpdateEndpoint          = "/v1/pages/{page_id}"
	APIDatabasesListEndpoint        = "/v1/databases"
	APIDatabasesRetrieveEndpoint    = "/v1/databases/{database_id}"
	APIDatabasesQueryEndpoint       = "/v1/databases/{database_id}/query"
	APISearchEndpoint               = "/v1/search"
)

const (
	DefaultNotionVersion = "2021-05-13"
)

type API struct {
	searchClient    SearchInterface
	usersClient     UsersInterface
	databasesClient DatabasesInterface
	pagesClient     PagesInterface
	blocksClient    BlocksInterface
}

func New(setters ...ClientOption) *API {
	options := defaultOptions

	for _, setter := range setters {
		setter(&options)
	}

	client := newHTTPClient(options)

	return &API{
		searchClient:    newSearchClient(client),
		usersClient:     newUsersClient(client),
		databasesClient: newDatabasesClient(client),
		pagesClient:     newPagesClient(client),
		blocksClient:    newBlocksClient(client),
	}
}

func (c *API) Users() UsersInterface {
	if c == nil {
		return nil
	}

	return c.usersClient
}

func (c *API) Databases() DatabasesInterface {
	if c == nil {
		return nil
	}

	return c.databasesClient
}

func (c *API) Pages() PagesInterface {
	if c == nil {
		return nil
	}

	return c.pagesClient
}

func (c *API) Blocks() BlocksInterface {
	if c == nil {
		return nil
	}

	return c.blocksClient
}

func (c *API) Search(ctx context.Context, params SearchParameters) (*SearchResponse, error) {
	if c == nil {
		return nil, ErrUnknown
	}

	return c.searchClient.Search(ctx, params)
}
