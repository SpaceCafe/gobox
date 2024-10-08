package types

// ServiceCreator is a generic interface for creating entities of type T.
type ServiceCreator[T any] interface {

	// Create creates a new entity in the given resource and returns the created entity.
	Create(resource Resource[T], options *ServiceOptions, entity *T) (createdEntity *T, err error)
}

// ServiceDeleter is a generic interface for deleting entities of type T.
type ServiceDeleter[T any] interface {

	// Delete removes an entity from the given resource.
	Delete(resource Resource[T], options *ServiceOptions, id string) error
}

// ServiceLister is a generic interface for listing entities of type T.
type ServiceLister[T any] interface {

	// List retrieves a list of entities from the given resource based on the provided options.
	List(resource Resource[T], options *ServiceOptions) (entities *[]T, meta *ServiceMeta, err error)
}

// ServiceReader is a generic interface for reading entities of type T.
type ServiceReader[T any] interface {

	// Read retrieves an entity from the given resource.
	Read(resource Resource[T], options *ServiceOptions, id string) (readEntity *T, err error)
}

// ServiceUpdater is a generic interface for updating entities of type T.
type ServiceUpdater[T any] interface {

	// Update modifies an existing entity in the given resource, either partially or fully.
	Update(resource Resource[T], options *ServiceOptions, id string, entity *T) (updatedEntity *T, err error)
}
