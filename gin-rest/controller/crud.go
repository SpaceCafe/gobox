package controller

// CRUD is a generic struct that combines all CRUD operations and capability handling.
// It embeds CapabilityFn, CreateFn, ReadFn, ListFn, UpdateFn, and DeleteFn to provide a complete set of handlers.
type CRUD[T any] struct {
	CapabilityFn
	CreateFn[T]
	ReadFn[T]
	ListFn[T]
	UpdateFn[T]
	DeleteFn[T]
}
