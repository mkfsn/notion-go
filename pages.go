package notion

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/mkfsn/notion-go/typed"
)

type Parent interface {
	isParent()
}

type baseParent struct {
	Type typed.ParentType `json:"type"`
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
	Object typed.ObjectType `json:"object"`
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

// FIXME: reduce the complexity
// nolint:gocyclo,gocognit,funlen
func (p *Page) UnmarshalJSON(data []byte) error {
	type Alias Page

	alias := struct {
		*Alias
		Parent     json.RawMessage            `json:"parent"`
		Properties map[string]json.RawMessage `json:"properties"`
	}{
		Alias: (*Alias)(p),
	}

	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}

	var baseParent baseParent

	if err := json.Unmarshal(alias.Parent, &baseParent); err != nil {
		return err
	}

	switch baseParent.Type {
	case typed.ParentTypeDatabase:
		var parent DatabaseParent

		if err := json.Unmarshal(alias.Parent, &parent); err != nil {
			return err
		}

		p.Parent = parent

	case typed.ParentTypePage:
		var parent PageParent

		if err := json.Unmarshal(alias.Parent, &parent); err != nil {
			return err
		}

		p.Parent = parent

	case typed.ParentTypeWorkspace:
		var parent WorkspaceParent

		if err := json.Unmarshal(alias.Parent, &parent); err != nil {
			return err
		}

		p.Parent = parent
	}

	p.Properties = make(map[string]PropertyValue)

	for name, value := range alias.Properties {
		var base basePropertyValue

		if err := json.Unmarshal(value, &base); err != nil {
			return err
		}

		switch base.Type {
		case typed.PropertyValueTypeRichText:
			var property RichTextPropertyValue

			if err := json.Unmarshal(value, &property); err != nil {
				return err
			}

			p.Properties[name] = property

		case typed.PropertyValueTypeNumber:
			var property NumberPropertyValue

			if err := json.Unmarshal(value, &property); err != nil {
				return err
			}

			p.Properties[name] = property

		case typed.PropertyValueTypeSelect:
			var property SelectPropertyValue

			if err := json.Unmarshal(value, &property); err != nil {
				return err
			}

			p.Properties[name] = property

		case typed.PropertyValueTypeMultiSelect:
			var property MultiSelectPropertyValue

			if err := json.Unmarshal(value, &property); err != nil {
				return err
			}

			p.Properties[name] = property

		case typed.PropertyValueTypeDate:
			var property DatePropertyValue

			if err := json.Unmarshal(value, &property); err != nil {
				return err
			}

			p.Properties[name] = property

		case typed.PropertyValueTypeFormula:
			var property FormulaPropertyValue

			if err := json.Unmarshal(value, &property); err != nil {
				return err
			}

			p.Properties[name] = property

		case typed.PropertyValueTypeRollup:
			var property RollupPropertyValue

			if err := json.Unmarshal(value, &property); err != nil {
				return err
			}

			p.Properties[name] = property

		case typed.PropertyValueTypeTitle:
			var property TitlePropertyValue

			if err := json.Unmarshal(value, &property); err != nil {
				return err
			}

			p.Properties[name] = property

		case typed.PropertyValueTypePeople:
			var property PeoplePropertyValue

			if err := json.Unmarshal(value, &property); err != nil {
				return err
			}

			p.Properties[name] = property

		case typed.PropertyValueTypeFiles:
			var property FilesPropertyValue

			if err := json.Unmarshal(value, &property); err != nil {
				return err
			}

			p.Properties[name] = property

		case typed.PropertyValueTypeCheckbox:
			var property CheckboxPropertyValue

			if err := json.Unmarshal(value, &property); err != nil {
				return err
			}

			p.Properties[name] = property

		case typed.PropertyValueTypeURL:
			var property URLPropertyValue

			if err := json.Unmarshal(value, &property); err != nil {
				return err
			}

			p.Properties[name] = property

		case typed.PropertyValueTypeEmail:
			var property EmailPropertyValue

			if err := json.Unmarshal(value, &property); err != nil {
				return err
			}

			p.Properties[name] = property

		case typed.PropertyValueTypePhoneNumber:
			var property PhoneNumberPropertyValue

			if err := json.Unmarshal(value, &property); err != nil {
				return err
			}

			p.Properties[name] = property

		case typed.PropertyValueTypeCreatedTime:
			var property CreatedTimePropertyValue

			if err := json.Unmarshal(value, &property); err != nil {
				return err
			}

			p.Properties[name] = property

		case typed.PropertyValueTypeCreatedBy:
			var property CreatedByPropertyValue

			if err := json.Unmarshal(value, &property); err != nil {
				return err
			}

			p.Properties[name] = property

		case typed.PropertyValueTypeLastEditedTime:
			var property LastEditedTimePropertyValue

			if err := json.Unmarshal(value, &property); err != nil {
				return err
			}

			p.Properties[name] = property

		case typed.PropertyValueTypeLastEditedBy:
			var property LastEditedByPropertyValue

			if err := json.Unmarshal(value, &property); err != nil {
				return err
			}

			p.Properties[name] = property
		}
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
	Type typed.PropertyValueType `json:"type,omitempty"`
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
		Title []json.RawMessage `json:"title"`
	}{
		Alias: (*Alias)(t),
	}

	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}

	t.Title = make([]RichText, 0, len(alias.Title))

	for _, value := range alias.Title {
		richText, err := newRichText(value)
		if err != nil {
			return err
		}

		t.Title = append(t.Title, richText)
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
		RichText []json.RawMessage `json:"rich_text"`
	}{
		Alias: (*Alias)(r),
	}

	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}

	r.RichText = make([]RichText, 0, len(alias.RichText))

	for _, value := range alias.RichText {
		richText, err := newRichText(value)
		if err != nil {
			return err
		}

		r.RichText = append(r.RichText, richText)
	}

	return nil
}

type NumberPropertyValue struct {
	basePropertyValue
	Number float64 `json:"number"`
}

type SelectPropertyValueOption struct {
	ID    string      `json:"id,omitempty"`
	Name  string      `json:"name"`
	Color typed.Color `json:"color,omitempty"`
}

type SelectPropertyValue struct {
	basePropertyValue
	Select SelectPropertyValueOption `json:"select"`
}

type MultiSelectPropertyValueOption struct {
	ID    string      `json:"id"`
	Name  string      `json:"name"`
	Color typed.Color `json:"color"`
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

func newFormulaValueType(data []byte) (FormulaValue, error) {
	var base baseFormulaValue

	if err := json.Unmarshal(data, &base); err != nil {
		return nil, err
	}

	switch base.Type {
	case typed.FormulaValueTypeString:
		var formulaValue StringFormulaValue

		if err := json.Unmarshal(data, &formulaValue); err != nil {
			return nil, err
		}

		return formulaValue, nil

	case typed.FormulaValueTypeNumber:
		var formulaValue NumberFormulaValue

		if err := json.Unmarshal(data, &formulaValue); err != nil {
			return nil, err
		}

		return formulaValue, nil

	case typed.FormulaValueTypeBoolean:
		var formulaValue BooleanFormulaValue

		if err := json.Unmarshal(data, &formulaValue); err != nil {
			return nil, err
		}

		return formulaValue, nil

	case typed.FormulaValueTypeDate:
		var formulaValue DateFormulaValue

		if err := json.Unmarshal(data, &formulaValue); err != nil {
			return nil, err
		}

		return formulaValue, nil
	}

	return nil, ErrUnknown
}

type baseFormulaValue struct {
	Type typed.FormulaValueType `json:"type"`
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
		Formula json.RawMessage `json:"formula"`
	}{
		Alias: (*Alias)(f),
	}

	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}

	formula, err := newFormulaValueType(alias.Formula)
	if err != nil {
		return err
	}

	f.Formula = formula

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
	endpoint := strings.Replace(APIPagesUpdateEndpoint, "{page_id}", params.PageID, 1)

	b, err := p.client.Request(ctx, http.MethodPatch, endpoint, params)
	if err != nil {
		return nil, err
	}

	var response PagesUpdateResponse

	if err := json.Unmarshal(b, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (p *pagesClient) Create(ctx context.Context, params PagesCreateParameters) (*PagesCreateResponse, error) {
	b, err := p.client.Request(ctx, http.MethodPost, APIPagesCreateEndpoint, params)
	if err != nil {
		return nil, err
	}

	var response PagesCreateResponse

	if err := json.Unmarshal(b, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
