package service

import (
	"github.com/spacecafe/gobox/gin-rest/types"
)

// CreateFn is a generic struct that implements the types.ServiceCreator interface.
type CreateFn[T any] struct{}

// Create is a method of CreateFn that handles the creation of a resource.
// It interacts with a repository (e.g. a database) to persist the entity and returns the created entity and any error encountered.
func (r *CreateFn[T]) Create(resource types.Resource[T], _ *types.ServiceOptions, entity *T) (createdEntity *T, err error) {
	result := resource.DB().Create(entity)
	return entity, result.Error
}
