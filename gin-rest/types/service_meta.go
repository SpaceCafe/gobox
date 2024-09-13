package types

import (
	"math"
)

// ServiceMeta contains metadata as a response for a view renderer.
type ServiceMeta struct {

	// Page is the current page number.
	Page int

	// PageSize is the number of items per page.
	PageSize int

	// Total is the total number of items available.
	Total int
}

// GetPage returns the current page number, ensuring it is at least 1.
func (r *ServiceMeta) GetPage() int {
	return Max(1, r.Page)
}

// GetPageSize returns the number of items per page, ensuring it is at least 1.
func (r *ServiceMeta) GetPageSize() int {
	return Max(1, r.PageSize)
}

// GetTotalPages calculates and returns the total number of pages, ensuring it is at least 1.
func (r *ServiceMeta) GetTotalPages() int {
	return Max(1, int(math.Ceil(float64(r.Total)/float64(r.GetPageSize()))))
}
