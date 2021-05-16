package typed

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
