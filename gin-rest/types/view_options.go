package types

// ViewOptions represents the options for viewing data, including filters, sorts, and MIME types.
type ViewOptions struct {

	// ServiceMeta contains metadata from the service about the passed entities.
	ServiceMeta

	// Filters is a pointer to a slice of FilterOption, which were used to filter the data.
	Filters *[]FilterOption

	// Sorts is a pointer to a slice of SortOption, which were used to sort the data.
	Sorts *[]SortOption
}
