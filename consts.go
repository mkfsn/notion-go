package notion

type BlockType string

const (
	BlockTypeParagraph        BlockType = "paragraph"
	BlockTypeHeading1         BlockType = "heading_1"
	BlockTypeHeading2         BlockType = "heading_2"
	BlockTypeHeading3         BlockType = "heading_3"
	BlockTypeBulletedListItem BlockType = "bulleted_list_item"
	BlockTypeNumberedListItem BlockType = "numbered_list_item"
	BlockTypeToDo             BlockType = "to_do"
	BlockTypeToggle           BlockType = "toggle"
	BlockTypeChildPage        BlockType = "child_page"
	BlockTypeUnsupported      BlockType = "unsupported"
)

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

type FormulaValueType string

const (
	FormulaValueTypeString  FormulaValueType = "string"
	FormulaValueTypeNumber  FormulaValueType = "number"
	FormulaValueTypeBoolean FormulaValueType = "boolean"
	FormulaValueTypeDate    FormulaValueType = "date"
)

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

type ObjectType string

const (
	ObjectTypeBlock    ObjectType = "block"
	ObjectTypePage     ObjectType = "page"
	ObjectTypeDatabase ObjectType = "database"
	ObjectTypeList     ObjectType = "list"
	ObjectTypeUser     ObjectType = "user"
)

type ParentType string

const (
	ParentTypeDatabase  ParentType = "database_id"
	ParentTypePage      ParentType = "page"
	ParentTypeWorkspace ParentType = "workspace"
)

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

type RichTextType string

const (
	RichTextTypeText     RichTextType = "text"
	RichTextTypeMention  RichTextType = "mention"
	RichTextTypeEquation RichTextType = "equation"
)

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

type UserType string

const (
	UserTypePerson UserType = "person"
	UserTypeBot    UserType = "bot"
)
