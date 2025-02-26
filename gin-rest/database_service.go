package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/spacecafe/gobox/gin-rest/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Ensure DatabaseService implements IService interface.
var _ types.IService = (*DatabaseService)(nil)

// DatabaseService provides methods to interact with the database using GORM and Gin.
type DatabaseService struct {
	// IResourceDatabase is an embedded field that represents the resource database interface.
	types.IResourceDatabase
}

// SetResource sets the resource for this service, ensuring it implements IResourceDatabase.
func (r *DatabaseService) SetResource(resource types.IResource) {
	databaseResource, ok := resource.(types.IResourceDatabase)
	if !ok {
		panic("resource must implement IResourceDatabase")
	}
	r.IResourceDatabase = databaseResource
}

// Create inserts a new entity into the database.
func (r *DatabaseService) Create(_ *gin.Context, entity any) (err error) {
	return r.Database().Create(entity).Error
}

// Read retrieves an entity from the database based on provided context and entity type.
func (r *DatabaseService) Read(ctx *gin.Context, entity any) (err error) {
	tx := r.Database().Preload(clause.Associations)
	if entity, ok := entity.(types.IModelReadable); ok {
		tx = tx.Select(entity.Readable(ctx))
	}
	return tx.First(entity).Error
}

// List retrieves a list of entities from the database based on provided context and entity type.
func (r *DatabaseService) List(ctx *gin.Context, entities any) (err error) {
	var total int64
	options := GetListOptions(ctx)
	err = options.Filter(r.Database()).Model(entities).Count(&total).Error
	if err != nil {
		return
	}
	ctx.Set(types.ContextDataTotal, int(total))

	tx := options.Prepare(r.Database())
	if entity, ok := ModelFactory(entities).(types.IModelReadable); ok {
		tx = tx.Select(entity.Readable(ctx))
	}
	return tx.Preload(clause.Associations).Find(entities).Error
}

// Update updates an existing entity in the database.
func (r *DatabaseService) Update(ctx *gin.Context, entity any) (err error) {
	tx := r.Database().Model(entity)
	if entity, ok := entity.(types.IModelUpdatable); ok {
		tx = tx.Select(entity.Updatable(ctx))
	}
	result := tx.Updates(entity)
	if result.Error == nil && result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

// Delete removes an entity from the database.
func (r *DatabaseService) Delete(_ *gin.Context, entity any) (err error) {
	result := r.Database().Model(entity).Delete(entity)
	if result.Error == nil && result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}
