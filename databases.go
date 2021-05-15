package notion

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

type Database struct {
	Object string `json:"object"`
	ID     string `json:"id"`

	CreatedTime    time.Time           `json:"created_time"`
	LastEditedTime time.Time           `json:"last_edited_time"`
	Title          []RichText          `json:"title"`
	Properties     map[string]Property `json:"properties"`
}

func (d Database) isSearchable() {}

func (d *Database) UnmarshalJSON(data []byte) error {
	type Alias Database

	alias := struct {
		*Alias
		Title      []json.RawMessage          `json:"title"`
		Properties map[string]json.RawMessage `json:"properties"`
	}{
		Alias: (*Alias)(d),
	}

	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}

	d.Title = make([]RichText, 0, len(alias.Title))
	d.Properties = make(map[string]Property)

	for _, title := range alias.Title {
		var base baseRichText

		if err := json.Unmarshal(title, &base); err != nil {
			return err
		}

		switch base.Type {
		case RichTextTypeText:
			var richText RichTextText

			if err := json.Unmarshal(title, &richText); err != nil {
				return err
			}

			d.Title = append(d.Title, richText)

		case RichTextTypeMention:
			var richText RichTextMention

			if err := json.Unmarshal(title, &richText); err != nil {
				return err
			}

			d.Title = append(d.Title, richText)

		case RichTextTypeEquation:
			var richText RichTextEquation

			if err := json.Unmarshal(title, &richText); err != nil {
				return err
			}

			d.Title = append(d.Title, richText)
		}
	}

	for name, value := range alias.Properties {
		var base baseProperty

		if err := json.Unmarshal(value, &base); err != nil {
			return err
		}

		switch base.Type {
		case PropertyTypeTitle:
			var property TitleProperty

			if err := json.Unmarshal(value, &property); err != nil {
				return err
			}

			d.Properties[name] = property

		case PropertyTypeRichText:
			var property RichTextProperty

			if err := json.Unmarshal(value, &property); err != nil {
				return err
			}

			d.Properties[name] = property

		case PropertyTypeNumber:
			var property NumberProperty

			if err := json.Unmarshal(value, &property); err != nil {
				return err
			}

			d.Properties[name] = property

		case PropertyTypeSelect:
			var property SelectProperty

			if err := json.Unmarshal(value, &property); err != nil {
				return err
			}

			d.Properties[name] = property

		case PropertyTypeMultiSelect:
			var property MultiSelectProperty

			if err := json.Unmarshal(value, &property); err != nil {
				return err
			}

			d.Properties[name] = property

		case PropertyTypeDate:
			var property DateProperty

			if err := json.Unmarshal(value, &property); err != nil {
				return err
			}

			d.Properties[name] = property

		case PropertyTypePeople:
			var property PeopleProperty

			if err := json.Unmarshal(value, &property); err != nil {
				return err
			}

			d.Properties[name] = property

		case PropertyTypeFile:
			var property FileProperty

			if err := json.Unmarshal(value, &property); err != nil {
				return err
			}

			d.Properties[name] = property

		case PropertyTypeCheckbox:
			var property CheckboxProperty

			if err := json.Unmarshal(value, &property); err != nil {
				return err
			}

			d.Properties[name] = property

		case PropertyTypeURL:
			var property URLProperty

			if err := json.Unmarshal(value, &property); err != nil {
				return err
			}

			d.Properties[name] = property

		case PropertyTypeEmail:
			var property EmailProperty

			if err := json.Unmarshal(value, &property); err != nil {
				return err
			}

			d.Properties[name] = property

		case PropertyTypePhoneNumber:
			var property PhoneNumberProperty

			if err := json.Unmarshal(value, &property); err != nil {
				return err
			}

			d.Properties[name] = property

		case PropertyTypeFormula:
			var property FormulaProperty

			if err := json.Unmarshal(value, &property); err != nil {
				return err
			}

			d.Properties[name] = property

		case PropertyTypeRelation:
			var property RelationProperty

			if err := json.Unmarshal(value, &property); err != nil {
				return err
			}

			d.Properties[name] = property

		case PropertyTypeRollup:
			var property RollupProperty

			if err := json.Unmarshal(value, &property); err != nil {
				return err
			}

			d.Properties[name] = property

		case PropertyTypeCreatedTime:
			var property CreatedTimeProperty

			if err := json.Unmarshal(value, &property); err != nil {
				return err
			}

			d.Properties[name] = property

		case PropertyTypeCreatedBy:
			var property CreatedByProperty

			if err := json.Unmarshal(value, &property); err != nil {
				return err
			}

			d.Properties[name] = property

		case PropertyTypeLastEditedTime:
			var property LastEditedTimeProperty

			if err := json.Unmarshal(value, &property); err != nil {
				return err
			}

			d.Properties[name] = property

		case PropertyTypeLastEditedBy:
			var property LastEditedByProperty

			if err := json.Unmarshal(value, &property); err != nil {
				return err
			}

			d.Properties[name] = property
		}
	}

	return nil
}

type Annotations struct {
	// Whether the text is **bolded**.
	Bold bool `json:"bold"`
	// Whether the text is _italicized_.
	Italic bool `json:"italic"`
	// Whether the text is ~~struck~~ through.
	Strikethrough bool `json:"strikethrough"`
	// Whether the text is __underlined__.
	Underline bool `json:"underline"`
	// Whether the text is `code style`.
	Code bool `json:"code"`
	// Color of the text.
	Color Color `json:"color"`
}

type RichText interface {
	isRichText()
}

func newRichText(data []byte) (RichText, error) {
	var base baseRichText

	if err := json.Unmarshal(data, &base); err != nil {
		return nil, err
	}

	switch base.Type {
	case RichTextTypeText:
		var richText RichTextText

		if err := json.Unmarshal(data, &richText); err != nil {
			return nil, err
		}

		return richText, nil

	case RichTextTypeMention:
		var richText RichTextMention

		if err := json.Unmarshal(data, &richText); err != nil {
			return nil, err
		}

		return richText, nil

	case RichTextTypeEquation:
		var richText RichTextEquation

		if err := json.Unmarshal(data, &richText); err != nil {
			return nil, err
		}

		return richText, nil
	}

	return nil, ErrUnknown
}

type RichTextType string

const (
	RichTextTypeText     RichTextType = "text"
	RichTextTypeMention  RichTextType = "mention"
	RichTextTypeEquation RichTextType = "equation"
)

type baseRichText struct {
	// The plain text without annotations.
	PlainText string `json:"plain_text"`
	// (Optional) The URL of any link or internal Notion mention in this text, if any.
	Href string `json:"href"`
	// Type of this rich text object.
	Type RichTextType `json:"type"`
	// All annotations that apply to this rich text.
	// Annotations include colors and bold/italics/underline/strikethrough.
	Annotations Annotations `json:"annotations"`
}

func (r baseRichText) isRichText() {}

type TextObject struct {
	Content string `json:"content"`
	Link    *struct {
		Type string `json:"type"`
		URL  string `json:"url"`
	} `json:"link"`
}

type RichTextText struct {
	baseRichText
	Text TextObject `json:"text"`
}

type Mention interface {
	isMention()
}

type baseMention struct {
	Type string `json:"type"`
}

func (b baseMention) isMention() {}

type UserMention struct {
	baseMention
	User User `json:"user"`
}

type PageMention struct {
	baseMention
	Page struct {
		ID string `json:"id"`
	} `json:"page"`
}

type DatabaseMention struct {
	baseMention
	Database struct {
		ID string `json:"id"`
	} `json:"database"`
}

type DateMention struct {
	baseMention
	Date DatePropertyValue `json:"date"`
}

type RichTextMention struct {
	baseRichText
	Mention Mention `json:"mention"`
}

type EquationObject struct {
	Expression string `json:"expression"`
}

type RichTextEquation struct {
	baseRichText
	Equation EquationObject `json:"equation"`
}

type Property interface {
	isProperty()
}

type PropertyType string

const (
	PropertyTypeTitle          PropertyType = "title"
	PropertyTypeRichText       PropertyType = "rich_text"
	PropertyTypeNumber         PropertyType = "number"
	PropertyTypeSelect         PropertyType = "select"
	PropertyTypeMultiSelect    PropertyType = "multi_select"
	PropertyTypeDate           PropertyType = "date"
	PropertyTypePeople         PropertyType = "people"
	PropertyTypeFile           PropertyType = "file"
	PropertyTypeCheckbox       PropertyType = "checkbox"
	PropertyTypeURL            PropertyType = "url"
	PropertyTypeEmail          PropertyType = "email"
	PropertyTypePhoneNumber    PropertyType = "phone_number"
	PropertyTypeFormula        PropertyType = "formula"
	PropertyTypeRelation       PropertyType = "relation"
	PropertyTypeRollup         PropertyType = "rollup"
	PropertyTypeCreatedTime    PropertyType = "created_time"
	PropertyTypeCreatedBy      PropertyType = "created_by"
	PropertyTypeLastEditedTime PropertyType = "last_edited_time"
	PropertyTypeLastEditedBy   PropertyType = "last_edited_by"
)

type baseProperty struct {
	// The ID of the property, usually a short string of random letters and symbols.
	// Some automatically generated property types have special human-readable IDs.
	// For example, all Title properties have an ID of "title".
	ID string `json:"id"`
	// Type that controls the behavior of the property
	Type PropertyType `json:"type"`
}

func (p baseProperty) isProperty() {}

type TitleProperty struct {
	baseProperty
	Title interface{} `json:"title"`
}

type RichTextProperty struct {
	baseProperty
	RichText interface{} `json:"rich_text"`
}

type NumberFormat string

const (
	NumberFormatNumber           NumberFormat = "number"
	NumberFormatNumberWithCommas NumberFormat = "number_with_commas"
	NumberFormatPercent          NumberFormat = "percent"
	NumberFormatDollar           NumberFormat = "dollar"
	NumberFormatEuro             NumberFormat = "euro"
	NumberFormatPound            NumberFormat = "pound"
	NumberFormatYen              NumberFormat = "yen"
	NumberFormatRuble            NumberFormat = "ruble"
	NumberFormatRupee            NumberFormat = "rupee"
	NumberFormatWon              NumberFormat = "won"
	NumberFormatYuan             NumberFormat = "yuan"
)

type NumberProperty struct {
	baseProperty
	Number struct {
		Format NumberFormat `json:"format"`
	} `json:"number"`
}

type SelectOption struct {
	Name  string `json:"name"`
	ID    string `json:"id"`
	Color Color  `json:"color"`
}

type MultiSelectOption struct {
	Name  string `json:"name"`
	ID    string `json:"id"`
	Color Color  `json:"color"`
}

type SelectProperty struct {
	baseProperty

	Select struct {
		Options []SelectOption `json:"options"`
	} `json:"select"`
}

type MultiSelectProperty struct {
	baseProperty
	MultiSelect struct {
		Options []MultiSelectOption `json:"options"`
	} `json:"multi_select"`
}

type DateProperty struct {
	baseProperty
	Date interface{} `json:"date"`
}

type PeopleProperty struct {
	baseProperty
	People interface{} `json:"people"`
}

type FileProperty struct {
	baseProperty
	File interface{} `json:"file"`
}

type CheckboxProperty struct {
	baseProperty
	Checkbox interface{} `json:"checkbox"`
}

type URLProperty struct {
	baseProperty
	URL interface{} `json:"url"`
}

type EmailProperty struct {
	baseProperty
	Email interface{} `json:"email"`
}

type PhoneNumberProperty struct {
	baseProperty
	PhoneNumber interface{} `json:"phone_number"`
}

type FormulaProperty struct {
	baseProperty
	Formula struct {
		Expression string `json:"expression"`
	} `json:"formula"`
}

type RelationProperty struct {
	baseProperty
	Relation struct {
		DatabaseID         string  `json:"database_id"`
		SyncedPropertyName *string `json:"synced_property_name"`
		SyncedPropertyID   *string `json:"synced_property_id"`
	} `json:"relation"`
}

type RollupFunction string

const (
	RollupFunctionCountAll          RollupFunction = "count_all"
	RollupFunctionCountValues       RollupFunction = "count_values"
	RollupFunctionCountUniqueValues RollupFunction = "count_unique_values"
	RollupFunctionCountEmpty        RollupFunction = "count_empty"
	RollupFunctionCountNotEmpty     RollupFunction = "count_not_empty"
	RollupFunctionPercentEmpty      RollupFunction = "percent_empty"
	RollupFunctionPercentNotEmpty   RollupFunction = "percent_not_empty"
	RollupFunctionSum               RollupFunction = "sum"
	RollupFunctionAverage           RollupFunction = "average"
	RollupFunctionMedian            RollupFunction = "median"
	RollupFunctionMin               RollupFunction = "min"
	RollupFunctionMax               RollupFunction = "max"
	RollupFunctionRange             RollupFunction = "range"
)

type RollupProperty struct {
	baseProperty
	Rollup struct {
		RelationPropertyName string         `json:"relation_property_name"`
		RelationPropertyID   string         `json:"relation_property_id"`
		RollupPropertyName   string         `json:"rollup_property_name"`
		RollupPropertyID     string         `json:"rollup_property_id"`
		Function             RollupFunction `json:"function"`
	} `json:"rollup"`
}

type CreatedTimeProperty struct {
	baseProperty
	CreatedTime interface{} `json:"created_time"`
}

type CreatedByProperty struct {
	baseProperty
	CreatedBy interface{} `json:"created_by"`
}

type LastEditedTimeProperty struct {
	baseProperty
	LastEditedTime interface{} `json:"last_edited_time"`
}

type LastEditedByProperty struct {
	baseProperty
	LastEditedBy interface{} `json:"last_edited_by"`
}

type DatabasesRetrieveParameters struct {
	DatabaseID string `json:"-"`
}

type DatabasesRetrieveResponse struct {
	Database
}

type DatabasesListParameters struct {
	PaginationParameters
}

type DatabasesListResponse struct {
	PaginatedList
	Results []Database `json:"results"`
}

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

type Sort struct {
	Property  string        `json:"property,omitempty"`
	Timestamp SortTimestamp `json:"timestamp,omitempty"`
	Direction SortDirection `json:"direction,omitempty"`
}

type Filter interface {
	isFilter()
}

type SinglePropertyFilter struct {
	Property string `json:"property"`
}

func (b SinglePropertyFilter) isFilter() {}

type TextFilter struct {
	Equals         *string `json:"equals,omitempty"`
	DoesNotEqual   *string `json:"does_not_equal,omitempty"`
	Contains       *string `json:"contains,omitempty"`
	DoesNotContain *string `json:"does_not_contain,omitempty"`
	StartsWith     *string `json:"starts_with,omitempty"`
	EndsWith       *string `json:"ends_with,omitempty"`
	IsEmpty        *bool   `json:"is_empty,omitempty"`
	IsNotEmpty     *bool   `json:"is_not_empty,omitempty"`
}

type SingleTextFilter struct {
	SinglePropertyFilter
	TextFilter
}

type NumberFilter struct {
	Equals               *float64 `json:"equals,omitempty"`
	DoesNotEqual         *float64 `json:"does_not_equal,omitempty"`
	GreaterThan          *float64 `json:"greater_than,omitempty"`
	LessThan             *float64 `json:"less_than,omitempty"`
	GreaterThanOrEqualTo *float64 `json:"greater_than_or_equal_to,omitempty"`
	LessThanOrEqualTo    *float64 `json:"less_than_or_equal_to,omitempty"`
	IsEmpty              *bool    `json:"is_empty,omitempty"`
	IsNotEmpty           *bool    `json:"is_not_empty,omitempty"`
}

type SingleNumberFilter struct {
	SinglePropertyFilter
	NumberFilter
}

type CheckboxFilter struct {
	Equals       *bool `json:"equals,omitempty"`
	DoesNotEqual *bool `json:"does_not_equal,omitempty"`
}

type CheckboxFilterWithProperty struct {
	SinglePropertyFilter
	CheckboxFilter
}

type SelectFilter struct {
	Equals       *string `json:"equals,omitempty"`
	DoesNotEqual *string `json:"does_not_equal,omitempty"`
	IsEmpty      *bool   `json:"is_empty,omitempty"`
	IsNotEmpty   *bool   `json:"is_not_empty,omitempty"`
}

type SingleSelectFilter struct {
	SinglePropertyFilter
	SelectFilter
}

type MultiSelectFilter struct {
	Contains       *string `json:"contains,omitempty"`
	DoesNotContain *string `json:"does_not_contain,omitempty"`
	IsEmpty        *bool   `json:"is_empty,omitempty"`
	IsNotEmpty     *bool   `json:"is_not_empty,omitempty"`
}

type SingleMultiSelectFilter struct {
	SinglePropertyFilter
	MultiSelectFilter
}

type DateFilter struct {
	Equals     *string                `json:"equals,omitempty"`
	Before     *string                `json:"before,omitempty"`
	After      *string                `json:"after,omitempty"`
	OnOrBefore *string                `json:"on_or_before,omitempty"`
	IsEmpty    *bool                  `json:"is_empty,omitempty"`
	IsNotEmpty *bool                  `json:"is_not_empty,omitempty"`
	OnOrAfter  *string                `json:"on_or_after,omitempty"`
	PastWeek   map[string]interface{} `json:"past_week,omitempty"`
	PastMonth  map[string]interface{} `json:"past_month,omitempty"`
	PastYear   map[string]interface{} `json:"past_year,omitempty"`
	NextWeek   map[string]interface{} `json:"next_week,omitempty"`
	NextMonth  map[string]interface{} `json:"next_month,omitempty"`
	NextYear   map[string]interface{} `json:"next_year,omitempty"`
}

type SingleDateFilter struct {
	SinglePropertyFilter
	DateFilter
}

type PeopleFilter struct {
	Contains       *string `json:"contains,omitempty"`
	DoesNotContain *string `json:"does_not_contain,omitempty"`
	IsEmpty        *bool   `json:"is_empty,omitempty"`
	IsNotEmpty     *bool   `json:"is_not_empty,omitempty"`
}

type SinglePeopleFilter struct {
	SinglePropertyFilter
	PeopleFilter
}

type FilesFilter struct {
	IsEmpty    *bool `json:"is_empty,omitempty"`
	IsNotEmpty *bool `json:"is_not_empty,omitempty"`
}

type SingleFilesFilter struct {
	SinglePropertyFilter
	FilesFilter
}

type RelationFilter struct {
	Contains       *string `json:"contains,omitempty"`
	DoesNotContain *string `json:"does_not_contain,omitempty"`
	IsEmpty        *bool   `json:"is_empty,omitempty"`
	IsNotEmpty     *bool   `json:"is_not_empty,omitempty"`
}

type SingleRelationFilter struct {
	SinglePropertyFilter
	RelationFilter
}

type FormulaFilter struct {
	Text     *TextFilter     `json:"text,omitempty"`
	Checkbox *CheckboxFilter `json:"checkbox,omitempty"`
	Number   *NumberFilter   `json:"number,omitempty"`
	Date     *DateFilter     `json:"date,omitempty"`
}

type SingleFormulaFilter struct {
	SinglePropertyFilter
	FormulaFilter
}

type CompoundFilter struct {
	Or  []Filter `json:"or,omitempty"`
	And []Filter `json:"and,omitempty"`
}

func (c CompoundFilter) isFilter() {}

type DatabasesQueryParameters struct {
	// Identifier for a Notion database.
	DatabaseID string `json:"-"`
	// When supplied, limits which pages are returned based on the
	// [filter conditions](https://developers.notion.com/reference-link/post-database-query-filter).
	Filter Filter `json:"filter"`
	// When supplied, orders the results based on the provided
	// [sort criteria](https://developers.notion.com/reference-link/post-database-query-sort).
	Sorts []Sort `json:"sorts"`
	// When supplied, returns a page of results starting after the cursor provided.
	// If not supplied, this endpoint will return the first page of results.
	StartCursor string `json:"-" url:"start_cursor,omitempty"`
	// The number of items from the full list desired in the response. Maximum: 100
	PageSize int32 `json:"-" url:"page_size,omitempty"`
}

type DatabasesQueryResponse struct {
	PaginatedList
	Results []Page `json:"results"`
}

type DatabasesInterface interface {
	Retrieve(ctx context.Context, params DatabasesRetrieveParameters) (*DatabasesRetrieveResponse, error)
	List(ctx context.Context, params DatabasesListParameters) (*DatabasesListResponse, error)
	Query(ctx context.Context, params DatabasesQueryParameters) (*DatabasesQueryResponse, error)
}

type databasesClient struct {
	client client
}

func newDatabasesClient(client client) *databasesClient {
	return &databasesClient{
		client: client,
	}
}

func (d *databasesClient) Retrieve(ctx context.Context, params DatabasesRetrieveParameters) (*DatabasesRetrieveResponse, error) {
	endpoint := strings.Replace(APIDatabasesRetrieveEndpoint, "{database_id}", params.DatabaseID, 1)

	b, err := d.client.Request(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	var response DatabasesRetrieveResponse

	if err := json.Unmarshal(b, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (d *databasesClient) List(ctx context.Context, params DatabasesListParameters) (*DatabasesListResponse, error) {
	b, err := d.client.Request(ctx, http.MethodGet, APIDatabasesListEndpoint, params)
	if err != nil {
		return nil, err
	}

	var response DatabasesListResponse

	if err := json.Unmarshal(b, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (d *databasesClient) Query(ctx context.Context, params DatabasesQueryParameters) (*DatabasesQueryResponse, error) {
	return nil, ErrUnimplemented
}
