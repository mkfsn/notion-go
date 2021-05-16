package typed

type FormulaValueType string

const (
	FormulaValueTypeString  FormulaValueType = "string"
	FormulaValueTypeNumber  FormulaValueType = "number"
	FormulaValueTypeBoolean FormulaValueType = "boolean"
	FormulaValueTypeDate    FormulaValueType = "date"
)
