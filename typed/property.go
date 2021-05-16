package typed

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

type PropertyValueType string

const (
	PropertyValueTypeRichText       PropertyValueType = "rich_text"
	PropertyValueTypeNumber         PropertyValueType = "number"
	PropertyValueTypeSelect         PropertyValueType = "select"
	PropertyValueTypeMultiSelect    PropertyValueType = "multi_select"
	PropertyValueTypeDate           PropertyValueType = "date"
	PropertyValueTypeFormula        PropertyValueType = "formula"
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
