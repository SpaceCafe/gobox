package service

import (
	"github.com/spacecafe/gobox/gin-rest/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// UpdateFn is a generic struct that implements the types.ServiceUpdater interface.
type UpdateFn[T any] struct{}

// Update is a method of UpdateFn that handles the updating of a resource.
// It interacts with a repository (e.g. a database) to update the entity either partially
// or fully based on the provided flag and returns the updated entity and any error encountered.
func (r *UpdateFn[T]) Update(resource types.Resource[T], options *types.ServiceOptions, entity *T) (updatedEntity *T, err error) {
	var result *gorm.DB

	if err = resource.DB().First(entity).Error; err != nil {
		return
	}

	if options.PartialUpdate {
		result = resource.DB().Clauses(clause.Returning{}).Updates(entity)
	} else {
		result = resource.DB().Clauses(clause.Returning{}).Save(entity)
	}

	if result.RowsAffected == 0 {
		return nil, types.ErrNotChanged
	}

	return entity, result.Error
}
