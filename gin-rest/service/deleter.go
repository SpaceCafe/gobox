package service

import (
	"github.com/spacecafe/gobox/gin-rest/types"
)

// DeleteFn is a generic struct that implements the types.ServiceDeleter interface.
type DeleteFn[T any] struct{}

// Delete is a method of DeleteFn that handles the deletion of a resource.
// It interacts with a repository (e.g. a database) to delete the entity and returns any error encountered.
func (r *DeleteFn[T]) Delete(resource types.Resource[T], _ *types.ServiceOptions, entity *T) error {
	result := resource.DB().Delete(entity)
	if result.RowsAffected == 0 {
		return types.ErrNotFound
	}
	return result.Error
}
