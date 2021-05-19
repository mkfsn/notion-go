package notion

type ObjectType string

const (
	ObjectTypeBlock    ObjectType = "block"
	ObjectTypePage     ObjectType = "page"
	ObjectTypeDatabase ObjectType = "database"
)

type PaginationParameters struct {
	// If supplied, this endpoint will return a page of results starting after the cursor provided.
	// If not supplied, this endpoint will return the first page of results.
	StartCursor string `json:"-" url:"start_cursor,omitempty"`
	// The number of items from the full list desired in the response. Maximum: 100
	PageSize int32 `json:"-" url:"page_size,omitempty"`
}

type PaginatedList struct {
	Object     string `json:"object"`
	HasMore    bool   `json:"has_more"`
	NextCursor string `json:"next_cursor"`
}

type Color string

const (
	DefaultColor          Color = "default"
	GrayColor             Color = "gray"
	BrownColor            Color = "brown"
	OrangeColor           Color = "orange"
	YellowColor           Color = "yellow"
	GreenColor            Color = "green"
	BlueColor             Color = "blue"
	PurpleColor           Color = "purple"
	PinkColor             Color = "pink"
	RedColor              Color = "red"
	GrayBackgroundColor   Color = "gray_background"
	BrownBackgroundColor  Color = "brown_background"
	OrangeBackgroundColor Color = "orange_background"
	YellowBackgroundColor Color = "yellow_background"
	GreenBackgroundColor  Color = "green_background"
	BlueBackgroundColor   Color = "blue_background"
	PurpleBackgroundColor Color = "purple_background"
	PinkBackgroundColor   Color = "pink_background"
	RedBackgroundColor    Color = "red_background"
)
