package notion

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/mkfsn/notion-go/rest"
)

type Database struct {
	Object ObjectType `json:"object"`
	ID     string     `json:"id"`

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
		Title      []richTextDecoder          `json:"title"`
		Properties map[string]propertyDecoder `json:"properties"`
	}{
		Alias: (*Alias)(d),
	}

	if err := json.Unmarshal(data, &alias); err != nil {
		return fmt.Errorf("failed to unmarshal Database: %w", err)
	}

	d.Title = make([]RichText, 0, len(alias.Title))

	for _, decoder := range alias.Title {
		d.Title = append(d.Title, decoder.RichText)
	}

	d.Properties = make(map[string]Property)

	for name, decoder := range alias.Properties {
		d.Properties[name] = decoder.Property
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

type BaseRichText struct {
	// The plain text without annotations.
	PlainText string `json:"plain_text,omitempty"`
	// (Optional) The URL of any link or internal Notion mention in this text, if any.
	Href string `json:"href,omitempty"`
	// Type of this rich text object.
	Type RichTextType `json:"type"`
	// All annotations that apply to this rich text.
	// Annotations include colors and bold/italics/underline/strikethrough.
	Annotations *Annotations `json:"annotations,omitempty"`
}

func (r BaseRichText) isRichText() {}

type Link struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}

type TextObject struct {
	Content string `json:"content"`
	Link    *Link  `json:"link,omitempty"`
}

type RichTextText struct {
	BaseRichText
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
	BaseRichText
	Mention Mention `json:"mention"`
}

type EquationObject struct {
	Expression string `json:"expression"`
}

type RichTextEquation struct {
	BaseRichText
	Equation EquationObject `json:"equation"`
}

type Property interface {
	isProperty()
}

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

type NumberPropertyOption struct {
	Format NumberFormat `json:"format"`
}

type NumberProperty struct {
	baseProperty
	Number NumberPropertyOption `json:"number"`
}

type SelectOption struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color Color  `json:"color"`
}

type MultiSelectOption struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color Color  `json:"color"`
}

type SelectPropertyOption struct {
	Options []SelectOption `json:"options"`
}

type SelectProperty struct {
	baseProperty
	Select SelectPropertyOption `json:"select"`
}

type MultiSelectPropertyOption struct {
	Options []MultiSelectOption `json:"options"`
}

type MultiSelectProperty struct {
	baseProperty
	MultiSelect MultiSelectPropertyOption `json:"multi_select"`
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

type Formula struct {
	Expression string `json:"expression"`
}

type FormulaProperty struct {
	baseProperty
	Formula Formula `json:"formula"`
}

type Relation struct {
	DatabaseID         string  `json:"database_id"`
	SyncedPropertyName *string `json:"synced_property_name"`
	SyncedPropertyID   *string `json:"synced_property_id"`
}

type RelationProperty struct {
	baseProperty
	Relation Relation `json:"relation"`
}

type RollupPropertyOption struct {
	RelationPropertyName string         `json:"relation_property_name"`
	RelationPropertyID   string         `json:"relation_property_id"`
	RollupPropertyName   string         `json:"rollup_property_name"`
	RollupPropertyID     string         `json:"rollup_property_id"`
	Function             RollupFunction `json:"function"`
}

type RollupProperty struct {
	baseProperty
	Rollup RollupPropertyOption `json:"rollup"`
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
	DatabaseID string `json:"-" url:"-"`
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

// SingleTextFilter is a text filter condition applies to database properties of types "title", "rich_text", "url", "email", and "phone".
type SingleTextFilter struct {
	SinglePropertyFilter
	Text     *TextFilter `json:"text,omitempty"`
	RichText *TextFilter `json:"rich_text,omitempty"`
	URL      *TextFilter `json:"url,omitempty"`
	Email    *TextFilter `json:"email,omitempty"`
	Phone    *TextFilter `json:"phone,omitempty"`
}

type NumberFilter struct {
	Equals               *float64 `json:"equals,omitempty"`
	DoesNotEqual         *float64 `json:"does_not_equal,omitempty"`
	GreaterThan          *float64 `json:"greater_than,omitempty"`
	LessThan             *float64 `json:"less_than,omitempty"`
	GreaterThanOrEqualTo *float64 `json:"greater_than_or_equal_to,omitempty"`
	LessThanOrEqualTo    *float64 `json:"less_than_or_equal_to,omitempty"`
	IsEmpty              bool     `json:"is_empty,omitempty"`
	IsNotEmpty           bool     `json:"is_not_empty,omitempty"`
}

// SingleNumberFilter is a number filter condition applies to database properties of type "number".
type SingleNumberFilter struct {
	SinglePropertyFilter
	Number NumberFilter `json:"number"`
}

type CheckboxFilter struct {
	Equals       bool `json:"equals,omitempty"`
	DoesNotEqual bool `json:"does_not_equal,omitempty"`
}

// SingleCheckboxFilter is a checkbox filter condition applies to database properties of type "checkbox".
type SingleCheckboxFilter struct {
	SinglePropertyFilter
	Checkbox CheckboxFilter `json:"checkbox"`
}

type SelectFilter struct {
	Equals       *string `json:"equals,omitempty"`
	DoesNotEqual *string `json:"does_not_equal,omitempty"`
	IsEmpty      bool    `json:"is_empty,omitempty"`
	IsNotEmpty   bool    `json:"is_not_empty,omitempty"`
}

// SingleSelectFilter is a select filter condition applies to database properties of type "select".
type SingleSelectFilter struct {
	SinglePropertyFilter
	Select SelectFilter `json:"select"`
}

type MultiSelectFilter struct {
	Contains       *string `json:"contains,omitempty"`
	DoesNotContain *string `json:"does_not_contain,omitempty"`
	IsEmpty        bool    `json:"is_empty,omitempty"`
	IsNotEmpty     bool    `json:"is_not_empty,omitempty"`
}

// SingleMultiSelectFilter is a multi-select filter condition applies to database properties of type "multi_select".
type SingleMultiSelectFilter struct {
	SinglePropertyFilter
	MultiSelect MultiSelectFilter `json:"multi_select"`
}

type DateFilter struct {
	Equals     *string                `json:"equals,omitempty"`
	Before     *string                `json:"before,omitempty"`
	After      *string                `json:"after,omitempty"`
	OnOrBefore *string                `json:"on_or_before,omitempty"`
	IsEmpty    bool                   `json:"is_empty,omitempty"`
	IsNotEmpty bool                   `json:"is_not_empty,omitempty"`
	OnOrAfter  *string                `json:"on_or_after,omitempty"`
	PastWeek   map[string]interface{} `json:"past_week,omitempty"`
	PastMonth  map[string]interface{} `json:"past_month,omitempty"`
	PastYear   map[string]interface{} `json:"past_year,omitempty"`
	NextWeek   map[string]interface{} `json:"next_week,omitempty"`
	NextMonth  map[string]interface{} `json:"next_month,omitempty"`
	NextYear   map[string]interface{} `json:"next_year,omitempty"`
}

// SingleDateFilter is a date filter condition applies to database properties of types "date", "created_time", and "last_edited_time".
type SingleDateFilter struct {
	SinglePropertyFilter
	Date           *DateFilter `json:"date,omitempty"`
	CreatedTime    *DateFilter `json:"created_time,omitempty"`
	LastEditedTime *DateFilter `json:"last_edited_time,omitempty"`
}

type PeopleFilter struct {
	Contains       *string `json:"contains,omitempty"`
	DoesNotContain *string `json:"does_not_contain,omitempty"`
	IsEmpty        bool    `json:"is_empty,omitempty"`
	IsNotEmpty     bool    `json:"is_not_empty,omitempty"`
}

// SinglePeopleFilter is a people filter condition applies to database properties of types "people", "created_by", and "last_edited_by".
type SinglePeopleFilter struct {
	SinglePropertyFilter
	People       *PeopleFilter `json:"people,omitempty"`
	CreatedBy    *PeopleFilter `json:"created_by,omitempty"`
	LastEditedBy *PeopleFilter `json:"last_edited_by,omitempty"`
}

type FilesFilter struct {
	IsEmpty    bool `json:"is_empty,omitempty"`
	IsNotEmpty bool `json:"is_not_empty,omitempty"`
}

// SingleFilesFilter is a files filter condition applies to database properties of type "files".
type SingleFilesFilter struct {
	SinglePropertyFilter
	Files FilesFilter `json:"files"`
}

type RelationFilter struct {
	Contains       *string `json:"contains,omitempty"`
	DoesNotContain *string `json:"does_not_contain,omitempty"`
	IsEmpty        bool    `json:"is_empty,omitempty"`
	IsNotEmpty     bool    `json:"is_not_empty,omitempty"`
}

// SingleRelationFilter is a relation filter condition applies to database properties of type "relation".
type SingleRelationFilter struct {
	SinglePropertyFilter
	Relation RelationFilter `json:"relation"`
}

type FormulaFilter struct {
	Text     *TextFilter     `json:"text,omitempty"`
	Checkbox *CheckboxFilter `json:"checkbox,omitempty"`
	Number   *NumberFilter   `json:"number,omitempty"`
	Date     *DateFilter     `json:"date,omitempty"`
}

// SingleFormulaFilter is a formula filter condition applies to database properties of type "formula".
type SingleFormulaFilter struct {
	SinglePropertyFilter
	Formula FormulaFilter `json:"formula"`
}

type CompoundFilter struct {
	Or  []Filter `json:"or,omitempty"`
	And []Filter `json:"and,omitempty"`
}

func (c CompoundFilter) isFilter() {}

type DatabasesQueryParameters struct {
	PaginationParameters
	// Identifier for a Notion database.
	DatabaseID string `json:"-" url:"-"`
	// When supplied, limits which pages are returned based on the
	// [filter conditions](https://developers.com/reference-link/post-database-query-filter).
	Filter Filter `json:"filter,omitempty" url:"-"`
	// When supplied, orders the results based on the provided
	// [sort criteria](https://developers.com/reference-link/post-database-query-sort).
	Sorts []Sort `json:"sorts,omitempty" url:"-"`
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
	restClient rest.Interface
}

func newDatabasesClient(restClient rest.Interface) *databasesClient {
	return &databasesClient{
		restClient: restClient,
	}
}

func (d *databasesClient) Retrieve(ctx context.Context, params DatabasesRetrieveParameters) (*DatabasesRetrieveResponse, error) {
	var result DatabasesRetrieveResponse

	var failure HTTPError

	err := d.restClient.New().Get().
		Endpoint(strings.Replace(APIDatabasesRetrieveEndpoint, "{database_id}", params.DatabaseID, 1)).
		Receive(ctx, &result, &failure)

	return &result, err // nolint:wrapcheck
}

func (d *databasesClient) List(ctx context.Context, params DatabasesListParameters) (*DatabasesListResponse, error) {
	var result DatabasesListResponse

	var failure HTTPError

	err := d.restClient.New().Get().
		Endpoint(APIDatabasesListEndpoint).
		QueryStruct(params).
		Receive(ctx, &result, &failure)

	return &result, err // nolint:wrapcheck
}

func (d *databasesClient) Query(ctx context.Context, params DatabasesQueryParameters) (*DatabasesQueryResponse, error) {
	var result DatabasesQueryResponse

	var failure HTTPError

	err := d.restClient.New().Post().
		Endpoint(strings.Replace(APIDatabasesQueryEndpoint, "{database_id}", params.DatabaseID, 1)).
		QueryStruct(params).
		BodyJSON(params).
		Receive(ctx, &result, &failure)

	return &result, err // nolint:wrapcheck
}

type richTextDecoder struct {
	RichText
}

func (r *richTextDecoder) UnmarshalJSON(data []byte) error {
	var decoder struct {
		Type RichTextType `json:"type"`
	}

	if err := json.Unmarshal(data, &decoder); err != nil {
		return fmt.Errorf("failed to unmarshal RichText: %w", err)
	}

	switch decoder.Type {
	case RichTextTypeText:
		r.RichText = &RichTextText{}

	case RichTextTypeMention:
		r.RichText = &RichTextMention{}

	case RichTextTypeEquation:
		r.RichText = &RichTextEquation{}
	}

	return json.Unmarshal(data, &r.RichText)
}

type propertyDecoder struct {
	Property
}

// UnmarshalJSON implements json.Unmarshaler
// nolint: cyclop
func (p *propertyDecoder) UnmarshalJSON(data []byte) error {
	var decoder struct {
		Type PropertyType `json:"type"`
	}

	if err := json.Unmarshal(data, &decoder); err != nil {
		return fmt.Errorf("failed to unmarshal Property: %w", err)
	}

	switch decoder.Type {
	case PropertyTypeTitle:
		p.Property = &TitleProperty{}

	case PropertyTypeRichText:
		p.Property = &RichTextProperty{}

	case PropertyTypeNumber:
		p.Property = &NumberProperty{}

	case PropertyTypeSelect:
		p.Property = &SelectProperty{}

	case PropertyTypeMultiSelect:
		p.Property = &MultiSelectProperty{}

	case PropertyTypeDate:
		p.Property = &DateProperty{}

	case PropertyTypePeople:
		p.Property = &PeopleProperty{}

	case PropertyTypeFile:
		p.Property = &FileProperty{}

	case PropertyTypeCheckbox:
		p.Property = &CheckboxProperty{}

	case PropertyTypeURL:
		p.Property = &URLProperty{}

	case PropertyTypeEmail:
		p.Property = &EmailProperty{}

	case PropertyTypePhoneNumber:
		p.Property = &PhoneNumberProperty{}

	case PropertyTypeFormula:
		p.Property = &FormulaProperty{}

	case PropertyTypeRelation:
		p.Property = &RelationProperty{}

	case PropertyTypeRollup:
		p.Property = &RollupProperty{}

	case PropertyTypeCreatedTime:
		p.Property = &CreatedTimeProperty{}

	case PropertyTypeCreatedBy:
		p.Property = &CreatedByProperty{}

	case PropertyTypeLastEditedTime:
		p.Property = &LastEditedTimeProperty{}

	case PropertyTypeLastEditedBy:
		p.Property = &LastEditedByProperty{}
	}

	return json.Unmarshal(data, &p.Property)
}
