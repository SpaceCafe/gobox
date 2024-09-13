package types

// ServiceCreator is a generic interface for creating entities of type T.
type ServiceCreator[T any] interface {

	// Create creates a new entity in the given resource and returns the created entity.
	Create(resource Resource[T], entity *T) (createdEntity *T, err error)
}

// ServiceDeleter is a generic interface for deleting entities of type T.
type ServiceDeleter[T any] interface {

	// Delete removes an entity from the given resource by its ID.
	Delete(resource Resource[T], id string) error
}

// ServiceLister is a generic interface for listing entities of type T.
type ServiceLister[T any] interface {

	// List retrieves a list of entities from the given resource based on the provided options.
	List(resource Resource[T], options *ServiceOptions) (entities *[]T, meta *ServiceMeta, err error)
}

// ServiceReader is a generic interface for reading entities of type T.
type ServiceReader[T any] interface {

	// Read retrieves an entity from the given resource by its ID.
	Read(resource Resource[T], id string) (entity *T, err error)
}

// ServiceUpdater is a generic interface for updating entities of type T.
type ServiceUpdater[T any] interface {

	// Update modifies an existing entity in the given resource, either partially or fully, based on the provided ID.
	Update(resource Resource[T], partially bool, id string, entity *T) (updatedEntity *T, err error)
}
