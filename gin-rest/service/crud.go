package service

// CRUD is a generic struct that combines all CRUD operations.
// It embeds CreateFn, ReadFn, ListFn, UpdateFn, and DeleteFn to provide a complete set of handlers.
type CRUD[T any] struct {
	CreateFn[T]
	ReadFn[T]
	ListFn[T]
	UpdateFn[T]
	DeleteFn[T]
}
