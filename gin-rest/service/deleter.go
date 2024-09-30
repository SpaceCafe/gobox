package service

import (
	"github.com/spacecafe/gobox/gin-rest/types"
	"gorm.io/gorm/clause"
)

// DeleteFn is a generic struct that implements the types.ServiceDeleter interface.
type DeleteFn[T any] struct{}

// Delete is a method of DeleteFn that handles the deletion of a resource.
// It interacts with a repository (e.g. a database) to delete the entity by its ID and returns any error encountered.
func (r *DeleteFn[T]) Delete(resource types.Resource[T], _ *types.ServiceOptions, id string) error {
	var entity T
	result := resource.DB().Clauses(clause.Where{Exprs: []clause.Expression{clause.Eq{Column: clause.PrimaryColumn, Value: id}}}).Delete(entity)
	if result.RowsAffected == 0 {
		return types.ErrNotFound
	}
	return result.Error
}
