package service

import (
	"github.com/spacecafe/gobox/gin-rest/types"
)

// ReadFn is a generic struct that implements the types.ServiceReader interface.
type ReadFn[T any] struct{}

// Read is a method of ReadFn that handles the retrieval of a resource.
// It interacts with a repository (e.g. a database) to fetch the entity and returns the entity and any error encountered.
func (r *ReadFn[T]) Read(resource types.Resource[T], _ *types.ServiceOptions, entity *T) (readEntity *T, err error) {
	result := resource.DB().First(entity)
	return entity, result.Error
}
