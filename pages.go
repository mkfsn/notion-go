package notion

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/mkfsn/notion-go/rest"
)

type Parent interface {
	isParent()
}

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
	Object ObjectType `json:"object"`
	// Unique identifier of the page.
	ID string `json:"id"`
	// The page's parent
	Parent Parent `json:"parent"`
	// Property values of this page.
	Properties map[string]PropertyValue `json:"properties"`
	// Date and time when this page was created. Formatted as an ISO 8601 date time string.
	CreatedTime time.Time `json:"created_time"`
	// Date and time when this page was updated. Formatted as an ISO 8601 date time string.
	LastEditedTime time.Time `json:"last_edited_time"`
	// The archived status of the page.
	Archived bool `json:"archived"`
}

func (p *Page) UnmarshalJSON(data []byte) error {
	type Alias Page

	alias := struct {
		*Alias
		Parent     parentDecoder                   `json:"parent"`
		Properties map[string]propertyValueDecoder `json:"properties"`
	}{
		Alias: (*Alias)(p),
	}

	if err := json.Unmarshal(data, &alias); err != nil {
		return fmt.Errorf("failed to unmarshal Page: %w", err)
	}

	p.Parent = alias.Parent.Parent

	p.Properties = make(map[string]PropertyValue)

	for name, decoder := range alias.Properties {
		p.Properties[name] = decoder.PropertyValue
	}

	return nil
}

func (p Page) isSearchable() {}

type PropertyValue interface {
	isPropertyValue()
}

type basePropertyValue struct {
	// Underlying identifier for the property. This identifier is guaranteed to remain constant when the property name changes.
	// It may be a UUID, but is often a short random string.
	// The id may be used in place of name when creating or updating pages.
	ID string `json:"id,omitempty"`
	// Type of the property
	Type PropertyValueType `json:"type,omitempty"`
}

func (p basePropertyValue) isPropertyValue() {}

type TitlePropertyValue struct {
	basePropertyValue
	Title []RichText `json:"title"`
}

func (t *TitlePropertyValue) UnmarshalJSON(data []byte) error {
	type Alias TitlePropertyValue

	alias := struct {
		*Alias
		Title []richTextDecoder `json:"title"`
	}{
		Alias: (*Alias)(t),
	}

	if err := json.Unmarshal(data, &alias); err != nil {
		return fmt.Errorf("failed to unmarshal TitlePropertyValue: %w", err)
	}

	t.Title = make([]RichText, 0, len(alias.Title))

	for _, decoder := range alias.Title {
		t.Title = append(t.Title, decoder.RichText)
	}

	return nil
}

type RichTextPropertyValue struct {
	basePropertyValue
	RichText []RichText `json:"rich_text"`
}

func (r *RichTextPropertyValue) UnmarshalJSON(data []byte) error {
	type Alias RichTextPropertyValue

	alias := struct {
		*Alias
		RichText []richTextDecoder `json:"rich_text"`
	}{
		Alias: (*Alias)(r),
	}

	if err := json.Unmarshal(data, &alias); err != nil {
		return fmt.Errorf("failed to unmarshal RichTextPropertyValue: %w", err)
	}

	r.RichText = make([]RichText, 0, len(alias.RichText))

	for _, decoder := range alias.RichText {
		r.RichText = append(r.RichText, decoder.RichText)
	}

	return nil
}

type NumberPropertyValue struct {
	basePropertyValue
	Number float64 `json:"number"`
}

type SelectPropertyValueOption struct {
	ID    string `json:"id,omitempty"`
	Name  string `json:"name"`
	Color Color  `json:"color,omitempty"`
}

type SelectPropertyValue struct {
	basePropertyValue
	Select SelectPropertyValueOption `json:"select"`
}

type MultiSelectPropertyValueOption struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color Color  `json:"color"`
}

type MultiSelectPropertyValue struct {
	basePropertyValue
	MultiSelect []MultiSelectPropertyValueOption `json:"multi_select"`
}

type DatePropertyValue struct {
	basePropertyValue
	Date struct {
		Start string  `json:"start"`
		End   *string `json:"end"`
	} `json:"date"`
}

type FormulaValue interface {
	isFormulaValue()
}

type baseFormulaValue struct {
	Type FormulaValueType `json:"type"`
}

func (b baseFormulaValue) isFormulaValue() {}

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
	Formula FormulaValue `json:"formula"`
}

func (f *FormulaPropertyValue) UnmarshalJSON(data []byte) error {
	type Alias FormulaPropertyValue

	alias := struct {
		*Alias
		Formula formulaValueDecoder `json:"formula"`
	}{
		Alias: (*Alias)(f),
	}

	if err := json.Unmarshal(data, &alias); err != nil {
		return fmt.Errorf("failed to unmarshal FormulaPropertyValue: %w", err)
	}

	f.Formula = alias.Formula.FormulaValue

	return nil
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
	PageID     string                   `json:"-" url:"-"`
	Properties map[string]PropertyValue `json:"properties" url:"-"`
}

type PagesUpdateResponse struct {
	Page
}

type PagesCreateParameters struct {
	// A DatabaseParentInput or PageParentInput
	Parent ParentInput `json:"parent" url:"-"`
	// Property values of this page. The keys are the names or IDs of the property and the values are property values.
	Properties map[string]PropertyValue `json:"properties" url:"-"`
	// Page content for the new page as an array of block objects
	Children []Block `json:"children,omitempty" url:"-"`
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
	restClient rest.Interface
}

func newPagesClient(restClient rest.Interface) *pagesClient {
	return &pagesClient{
		restClient: restClient,
	}
}

func (p *pagesClient) Retrieve(ctx context.Context, params PagesRetrieveParameters) (*PagesRetrieveResponse, error) {
	var result PagesRetrieveResponse

	var failure HTTPError

	err := p.restClient.New().Get().
		Endpoint(strings.Replace(APIPagesRetrieveEndpoint, "{page_id}", params.PageID, 1)).
		QueryStruct(params).
		BodyJSON(nil).
		Receive(ctx, &result, &failure)

	return &result, err // nolint:wrapcheck
}

func (p *pagesClient) Update(ctx context.Context, params PagesUpdateParameters) (*PagesUpdateResponse, error) {
	var result PagesUpdateResponse

	var failure HTTPError

	err := p.restClient.New().Patch().
		Endpoint(strings.Replace(APIPagesUpdateEndpoint, "{page_id}", params.PageID, 1)).
		QueryStruct(params).
		BodyJSON(params).
		Receive(ctx, &result, &failure)

	return &result, err // nolint:wrapcheck
}

func (p *pagesClient) Create(ctx context.Context, params PagesCreateParameters) (*PagesCreateResponse, error) {
	var result PagesCreateResponse

	var failure HTTPError

	err := p.restClient.New().Post().
		Endpoint(APIPagesCreateEndpoint).
		QueryStruct(params).
		BodyJSON(params).
		Receive(ctx, &result, &failure)

	return &result, err // nolint:wrapcheck
}

type formulaValueDecoder struct {
	FormulaValue
}

func (f *formulaValueDecoder) UnmarshalJSON(data []byte) error {
	var decoder struct {
		Type FormulaValueType `json:"type"`
	}

	if err := json.Unmarshal(data, &decoder); err != nil {
		return fmt.Errorf("failed to unmarshal FormulaValue: %w", err)
	}

	switch decoder.Type {
	case FormulaValueTypeString:
		f.FormulaValue = &StringFormulaValue{}

	case FormulaValueTypeNumber:
		f.FormulaValue = &NumberFormulaValue{}

	case FormulaValueTypeBoolean:
		f.FormulaValue = &BooleanFormulaValue{}

	case FormulaValueTypeDate:
		f.FormulaValue = &DateFormulaValue{}
	}

	return json.Unmarshal(data, &f.FormulaValue)
}

type parentDecoder struct {
	Parent
}

func (p *parentDecoder) UnmarshalJSON(data []byte) error {
	var decoder struct {
		Type ParentType `json:"type"`
	}

	if err := json.Unmarshal(data, &decoder); err != nil {
		return fmt.Errorf("failed to unmarshal Parent: %w", err)
	}

	switch decoder.Type {
	case ParentTypeDatabase:
		p.Parent = &DatabaseParent{}

	case ParentTypePage:
		p.Parent = &PageParent{}

	case ParentTypeWorkspace:
		p.Parent = &WorkspaceParent{}
	}

	return json.Unmarshal(data, p.Parent)
}

type propertyValueDecoder struct {
	PropertyValue
}

// UnmarshalJSON implements json.Unmarshaler
// nolint: cyclop
func (p *propertyValueDecoder) UnmarshalJSON(data []byte) error {
	var decoder struct {
		Type PropertyValueType `json:"type,omitempty"`
	}

	if err := json.Unmarshal(data, &decoder); err != nil {
		return fmt.Errorf("failed to unmarshal PropertyValue: %w", err)
	}

	switch decoder.Type {
	case PropertyValueTypeRichText:
		p.PropertyValue = &RichTextPropertyValue{}

	case PropertyValueTypeNumber:
		p.PropertyValue = &NumberPropertyValue{}

	case PropertyValueTypeSelect:
		p.PropertyValue = &SelectPropertyValue{}

	case PropertyValueTypeMultiSelect:
		p.PropertyValue = &MultiSelectPropertyValue{}

	case PropertyValueTypeDate:
		p.PropertyValue = &DatePropertyValue{}

	case PropertyValueTypeFormula:
		p.PropertyValue = &FormulaPropertyValue{}

	case PropertyValueTypeRollup:
		p.PropertyValue = &RollupPropertyValue{}

	case PropertyValueTypeTitle:
		p.PropertyValue = &TitlePropertyValue{}

	case PropertyValueTypePeople:
		p.PropertyValue = &PeoplePropertyValue{}

	case PropertyValueTypeFiles:
		p.PropertyValue = &FilesPropertyValue{}

	case PropertyValueTypeCheckbox:
		p.PropertyValue = &CheckboxPropertyValue{}

	case PropertyValueTypeURL:
		p.PropertyValue = &URLPropertyValue{}

	case PropertyValueTypeEmail:
		p.PropertyValue = &EmailPropertyValue{}

	case PropertyValueTypePhoneNumber:
		p.PropertyValue = &PhoneNumberPropertyValue{}

	case PropertyValueTypeCreatedTime:
		p.PropertyValue = &CreatedTimePropertyValue{}

	case PropertyValueTypeCreatedBy:
		p.PropertyValue = &CreatedByPropertyValue{}

	case PropertyValueTypeLastEditedTime:
		p.PropertyValue = &LastEditedTimePropertyValue{}

	case PropertyValueTypeLastEditedBy:
		p.PropertyValue = &LastEditedByPropertyValue{}
	}

	return json.Unmarshal(data, p.PropertyValue)
}
