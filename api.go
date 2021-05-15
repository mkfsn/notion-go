package notion

const (
	APIBaseURL                        = "https://api.notion.com"
	APIUsersListEndpoint              = "/v1/users"
	APIUsersRetrieveEndpoint          = "/v1/users/{user_id}"
	APIBlocksRetrieveChildrenEndpoint = "/v1/blocks/{block_id}/children"
	APIBlocksAppendChildrenEndpoint   = "/v1/blocks/{block_id}/children"
	APIPagesCreateEndpoint            = "/v1/pages"
	APIPagesRetrieveEndpoint          = "/v1/pages/{page_id}"
	APIPagesUpdateEndpoint            = "/v1/pages/{page_id}"
	APIDatabasesListEndpoint          = "/v1/databases"
	APIDatabasesRetrieveEndpoint      = "/v1/databases/{database_id}"
	APIDatabasesQueryEndpoint         = "/v1/databases/{database_id}/query"
	APISearchEndpoint                 = "/v1/search"
)

const (
	DefaultNotionVersion = "2021-05-13"
)

type API struct {
	usersClient     *usersClient
	databasesClient *databasesClient
	pagesClient     *pagesClient
	blocksClient    *blocksClient
}

func New(setters ...ClientOption) *API {
	options := defaultOptions
	for _, setter := range setters {
		setter(&options)
	}

	client := newHTTPClient(options)

	return &API{
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
