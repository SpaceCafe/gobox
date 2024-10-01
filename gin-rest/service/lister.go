package service

import (
	"github.com/spacecafe/gobox/gin-rest/types"
	"gorm.io/gorm/clause"
)

// ListFn is a generic struct that implements the types.ServiceLister interface.
type ListFn[T any] struct{}

// List is a method of ListFn that handles the listing of resources.
// It interacts with a repository (e.g. a database) to retrieve entities based on the provided options
// and returns the entities, metadata, and any error encountered.
func (r *ListFn[T]) List(resource types.Resource[T], options *types.ServiceOptions) (entities *[]T, meta *types.ServiceMeta, err error) {
	entities = &[]T{}
	result := ApplyOptions(resource.DB(), options).Preload(clause.Associations).Find(entities)
	return entities, &types.ServiceMeta{
		Page:     options.Page,
		PageSize: options.PageSize,
		Total:    int(result.RowsAffected),
	}, result.Error
}
