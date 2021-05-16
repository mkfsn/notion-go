package typed

type SortTimestamp string

const (
	SortTimestampByCreatedTime    SortTimestamp = "created_time"
	SortTimestampByLastEditedTime SortTimestamp = "last_edited_time"
)

type SortDirection string

const (
	SortDirectionAscending  SortDirection = "ascending"
	SortDirectionDescending SortDirection = "descending"
)
