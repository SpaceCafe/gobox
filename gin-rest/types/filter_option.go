package types

// FilterOption represents a single filter criterion for querying.
type FilterOption struct {

	// Field is the name of the field to filter by.
	Field string

	// Value is the value to filter the field by.
	Value string

	// Operator is the operator to use for filtering (e.g., equals, less than).
	Operator FilterOperator
}

// FilterOperator represents the type of operation to perform for filtering.
type FilterOperator string

const (
	// Equals checks if the field value is equal to the specified value.
	Equals FilterOperator = "eq"

	// NotEquals checks if the field value is not equal to the specified value.
	NotEquals FilterOperator = "ne"

	// GreaterThan checks if the field value is greater than the specified value.
	GreaterThan FilterOperator = "gt"

	// GreaterThanOrEqual checks if the field value is greater than or equal to the specified value.
	GreaterThanOrEqual FilterOperator = "gte"

	// LessThan checks if the field value is less than the specified value.
	LessThan FilterOperator = "lt"

	// LessThanOrEqual checks if the field value is less than or equal to the specified value.
	LessThanOrEqual FilterOperator = "lte"

	// Like checks if the field value matches the specified pattern.
	Like FilterOperator = "like"

	// Fuzzy checks if the field value approximately matches the specified value.
	Fuzzy FilterOperator = "fuzzy"
)

var (
	// FilterOperators is a map that holds a set of predefined filter operators.
	// This map is used to check if a given operator is valid.
	//nolint:gochecknoglobals // Maintain a set of predefined operators that are used throughout the application.
	FilterOperators = map[FilterOperator]struct{}{
		Equals:             {},
		NotEquals:          {},
		GreaterThan:        {},
		GreaterThanOrEqual: {},
		LessThan:           {},
		LessThanOrEqual:    {},
		Like:               {},
		Fuzzy:              {},
	}
)
