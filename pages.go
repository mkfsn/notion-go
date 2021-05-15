package notion

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

type Parent interface {
	isParent()
}

type ParentType string

const (
	ParentTypeDatabase  ParentType = "database_id"
	ParentTypePage      ParentType = "page"
	ParentTypeWorkspace ParentType = "workspace"
)

type baseParent struct {
	Type ParentType `json:"type"`
}

func (b baseParent) isParent() {}

type DatabaseParent struct {
	baseParent
	DatabaseID string `json:"database_id"`
}

type PageParent struct {
	baseParent
	PageID string `json:"page_id"`
}

type WorkspaceParent struct {
	baseParent
}

type ParentInput interface {
	isParentInput()
}

type baseParentInput struct{}

func (b baseParentInput) isParentInput() {}

type DatabaseParentInput struct {
	baseParentInput
	DatabaseID string `json:"database_id"`
}

type PageParentInput struct {
	baseParentInput
	PageID string `json:"page_id"`
}

type Page struct {
	// Always "page".
	Object string `json:"object"`
	// Unique identifier of the page.
	ID string `json:"string"`
	// The page's parent
	Parent Parent `json:"parent"`
	// Property values of this page.
	Properties map[string]PropertyValue `json:"properties"`

	// TODO: check if the following fields are valid
	// Date and time when this page was created. Formatted as an ISO 8601 date time string.
	CreatedTime time.Time `json:"created_time"`
	// Date and time when this page was updated. Formatted as an ISO 8601 date time string.
	LastEditedTime time.Time `json:"last_edited_time"`
	// The archived status of the page.
	Archived bool `json:"archived"`
}

type PropertyValue interface {
	isPropertyValue()
}

type PropertyValueType string

const (
	PropertyValueTypeRichText       PropertyValueType = "rich_text"
	PropertyValueTypeNumber         PropertyValueType = "number"
	PropertyValueTypeSelect         PropertyValueType = "select"
	PropertyValueTypeMultiSelect    PropertyValueType = "multi_select"
	PropertyValueTypeDate           PropertyValueType = "date"
	PropertyValueTypeFormula        PropertyValueType = "formula"
	PropertyValueTypeRelation       PropertyValueType = "relation"
	PropertyValueTypeRollup         PropertyValueType = "rollup"
	PropertyValueTypeTitle          PropertyValueType = "title"
	PropertyValueTypePeople         PropertyValueType = "people"
	PropertyValueTypeFiles          PropertyValueType = "files"
	PropertyValueTypeCheckbox       PropertyValueType = "checkbox"
	PropertyValueTypeURL            PropertyValueType = "url"
	PropertyValueTypeEmail          PropertyValueType = "email"
	PropertyValueTypePhoneNumber    PropertyValueType = "phone_number"
	PropertyValueTypeCreatedTime    PropertyValueType = "created_time"
	PropertyValueTypeCreatedBy      PropertyValueType = "created_by"
	PropertyValueTypeLastEditedTime PropertyValueType = "last_edited_time"
	PropertyValueTypeLastEditedBy   PropertyValueType = "last_edited_by"
)

type basePropertyValue struct {
	// Underlying identifier for the property. This identifier is guaranteed to remain constant when the property name changes.
	// It may be a UUID, but is often a short random string.
	// The id may be used in place of name when creating or updating pages.
	ID string `json:"id"`
	// Type of the property
	Type PropertyValueType `json:"type"`
}

func (p basePropertyValue) isPropertyValue() {}

type TitlePropertyValue struct {
	basePropertyValue
	Title []RichText `json:"title"`
}

type RichTextPropertyValue struct {
	basePropertyValue
	RichText []RichText `json:"rich_text"`
}

type NumberPropertyValue struct {
	basePropertyValue
	Number float64 `json:"number"`
}

type SelectPropertyValue struct {
	basePropertyValue
	Select struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Color Color  `json:"color"`
	} `json:"select"`
}

type MultiSelectPropertyValue struct {
	basePropertyValue
	MultiSelect struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Color Color  `json:"color"`
	} `json:"multi_select"`
}

type DatePropertyValue struct {
	basePropertyValue
	Date struct {
		Start string  `json:"start"`
		End   *string `json:"end"`
	} `json:"date"`
}

type FormulaValueType interface {
	isFormulaValueType()
}

type baseFormulaValue struct {
	Type string `json:"type"`
}

func (b baseFormulaValue) isFormulaValueType() {}

type StringFormulaValue struct {
	baseFormulaValue
	String *string `json:"string"`
}

type NumberFormulaValue struct {
	baseFormulaValue
	Number *float64 `json:"number"`
}

type BooleanFormulaValue struct {
	baseFormulaValue
	Boolean bool `json:"boolean"`
}

type DateFormulaValue struct {
	baseFormulaValue
	Date DatePropertyValue `json:"date"`
}

type FormulaPropertyValue struct {
	basePropertyValue
	Formula FormulaValueType `json:"formula"`
}

type RollupValueType interface {
	isRollupValueType()
}

type baseRollupValueType struct {
	Type string `json:"type"`
}

func (b baseRollupValueType) isRollupValueType() {}

type NumberRollupValue struct {
	baseRollupValueType
	Number float64 `json:"number"`
}

type DateRollupValue struct {
	baseRollupValueType
	Date DatePropertyValue `json:"date"`
}

type ArrayRollupValue struct {
	baseRollupValueType
	Array []interface{} `json:"array"`
}

type RollupPropertyValue struct {
	basePropertyValue
	Rollup RollupValueType `json:"rollup"`
}

type PeoplePropertyValue struct {
	basePropertyValue
	People []User `json:"people"`
}

type FilesPropertyValue struct {
	basePropertyValue
	Files []struct {
		Name string `json:"name"`
	} `json:"files"`
}

type CheckboxPropertyValue struct {
	basePropertyValue
	Checkbox bool `json:"checkbox"`
}

type URLPropertyValue struct {
	basePropertyValue
	URL string `json:"url"`
}

type EmailPropertyValue struct {
	basePropertyValue
	Email string `json:"email"`
}

type PhoneNumberPropertyValue struct {
	basePropertyValue
	PhoneNumber string `json:"phone_number"`
}

type CreatedTimePropertyValue struct {
	basePropertyValue
	CreatedTime time.Time `json:"created_time"`
}

type CreatedByPropertyValue struct {
	basePropertyValue
	CreatedBy User `json:"created_by"`
}

type LastEditedTimePropertyValue struct {
	basePropertyValue
	LastEditedTime time.Time `json:"last_edited_time"`
}

type LastEditedByPropertyValue struct {
	basePropertyValue
	LastEditedBy User `json:"last_edited_by"`
}

type PagesRetrieveParameters struct {
	PageID string
}

type PagesRetrieveResponse struct {
	Page
}

type PagesUpdateParameters struct {
	PageID     string
	Properties map[string]PropertyValue `json:"properties"`
}

type PagesUpdateResponse struct {
	Page
}

type PagesCreateParameters struct {
	// A DatabaseParentInput or PageParentInput
	Parent ParentInput `json:"parent"`
	// Property values of this page. The keys are the names or IDs of the property and the values are property values.
	Properties map[string]PropertyValue `json:"properties"`
	// Page content for the new page as an array of block objects
	Children []Block `json:"children"`
}

type PagesCreateResponse struct {
	Page
}

type PagesInterface interface {
	Retrieve(ctx context.Context, params PagesRetrieveParameters) (*PagesRetrieveResponse, error)
	Update(ctx context.Context, params PagesUpdateParameters) (*PagesUpdateResponse, error)
	Create(ctx context.Context, params PagesCreateParameters) (*PagesCreateResponse, error)
}

type pagesClient struct {
	client client
}

func newPagesClient(client client) *pagesClient {
	return &pagesClient{
		client: client,
	}
}

func (p *pagesClient) Retrieve(ctx context.Context, params PagesRetrieveParameters) (*PagesRetrieveResponse, error) {
	endpoint := strings.Replace(APIPagesRetrieveEndpoint, "{page_id}", params.PageID, 1)

	b, err := p.client.Request(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	var response PagesRetrieveResponse

	if err := json.Unmarshal(b, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (p *pagesClient) Update(ctx context.Context, params PagesUpdateParameters) (*PagesUpdateResponse, error) {
	return nil, ErrUnimplemented
}

func (p *pagesClient) Create(ctx context.Context, params PagesCreateParameters) (*PagesCreateResponse, error) {
	return nil, ErrUnimplemented
}
