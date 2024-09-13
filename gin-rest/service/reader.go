package service

import (
	"github.com/spacecafe/gobox/gin-rest/types"
	"gorm.io/gorm/clause"
)

// ReadFn is a generic struct that implements the types.ServiceReader interface.
type ReadFn[T any] struct{}

// Read is a method of ReadFn that handles the retrieval of a resource.
// It interacts with a repository (e.g. a database) to fetch the entity by its ID and returns the entity and any error encountered.
func (r *ReadFn[T]) Read(resource types.Resource[T], id string) (entity *T, err error) {
	entity = new(T)
	result := resource.DB().Clauses(clause.Where{Exprs: []clause.Expression{clause.Eq{Column: clause.PrimaryColumn, Value: id}}}).First(entity)
	return entity, result.Error
}
