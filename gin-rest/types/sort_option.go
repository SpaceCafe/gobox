package types

// SortOption represents a single sort criterion for querying.
type SortOption struct {

	// Field is the name of the field to sort by.
	Field string

	// Descending indicates whether the sorting should be in descending order.
	Descending bool
}
