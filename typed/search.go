package typed

type SearchFilterValue string

const (
	SearchFilterValuePage     SearchFilterValue = "page"
	SearchFilterValueDatabase SearchFilterValue = "database"
)

type SearchFilterProperty string

const (
	SearchFilterPropertyObject SearchFilterProperty = "object"
)

type SearchSortDirection string

const (
	SearchSortDirectionAscending  SearchSortDirection = "ascending"
	SearchSortDirectionDescending SearchSortDirection = " descending"
)

type SearchSortTimestamp string

const (
	SearchSortTimestampLastEditedTime SearchSortTimestamp = "last_edited_time"
)
