package notion

import (
	"context"
	"net/http"

	"github.com/mkfsn/notion-go/rest"
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
	DefaultUserAgent     = "mkfsn/notion-go"
)

var (
	defaultSettings = apiSettings{
		baseURL:       APIBaseURL,
		notionVersion: DefaultNotionVersion,
		userAgent:     DefaultUserAgent,
		httpClient:    http.DefaultClient,
	}
)

type API struct {
	searchClient    SearchInterface
	usersClient     UsersInterface
	databasesClient DatabasesInterface
	pagesClient     PagesInterface
	blocksClient    BlocksInterface
}

func New(authToken string, setters ...APISetting) *API {
	settings := defaultSettings

	for _, setter := range setters {
		setter(&settings)
	}

	restClient := rest.New().
		BearerToken(authToken).
		BaseURL(settings.baseURL).
		UserAgent(settings.userAgent).
		Client(settings.httpClient).
		Header("Notion-Version", settings.notionVersion)

	return &API{
		searchClient:    newSearchClient(restClient),
		usersClient:     newUsersClient(restClient),
		databasesClient: newDatabasesClient(restClient),
		pagesClient:     newPagesClient(restClient),
		blocksClient:    newBlocksClient(restClient),
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

type apiSettings struct {
	baseURL       string
	notionVersion string
	userAgent     string
	httpClient    *http.Client
}

type APISetting func(o *apiSettings)

func WithBaseURL(baseURL string) APISetting {
	return func(o *apiSettings) {
		o.baseURL = baseURL
	}
}

func WithNotionVersion(notionVersion string) APISetting {
	return func(o *apiSettings) {
		o.notionVersion = notionVersion
	}
}

func WithUserAgent(userAgent string) APISetting {
	return func(o *apiSettings) {
		o.userAgent = userAgent
	}
}

func WithHTTPClient(httpClient *http.Client) APISetting {
	return func(o *apiSettings) {
		o.httpClient = httpClient
	}
}
