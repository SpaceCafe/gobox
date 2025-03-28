package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/spacecafe/gobox/gin-rest/types"
	"gorm.io/gorm"
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
func (r *DatabaseService) Create(ctx *gin.Context, entity any) (err error) {
	stmt := r.Database().Statement
	if entity, ok := entity.(types.IModelCreateClause); ok {
		AddClauses(stmt, entity.CreateClause(ctx))
	}
	return stmt.Create(entity).Error
}

// Read retrieves an entity from the database based on provided context and entity type.
func (r *DatabaseService) Read(ctx *gin.Context, entity any) (err error) {
	stmt := r.Database().Statement
	if entity, ok := entity.(types.IModelReadable); ok {
		stmt.Select(entity.Readable(ctx))
	}
	if entity, ok := entity.(types.IModelReadClause); ok {
		AddClauses(stmt, entity.ReadClause(ctx))
	}
	return stmt.First(entity).Error
}

// List retrieves a list of entities from the database based on provided context and entity type.
func (r *DatabaseService) List(ctx *gin.Context, entities any) (err error) {
	var total int64
	entity := ModelFactory(entities)

	options := GetListOptions(ctx)
	filter := options.Filter()

	stmt := r.Database().Model(entities).Statement
	stmt.AddClause(filter)

	if entity, ok := entity.(types.IModelReadable); ok {
		stmt.Select(entity.Readable(ctx))
	}
	if entity, ok := entity.(types.IModelListClause); ok {
		AddClauses(stmt, entity.ListClause(ctx))
	}

	err = stmt.Count(&total).Error
	if err != nil {
		return
	}
	ctx.Set(types.ContextDataTotal, int(total))

	stmt.AddClause(options.Paginate())
	stmt.AddClause(options.Sort())
	return stmt.Find(entities).Error
}

// Update updates an existing entity in the database.
func (r *DatabaseService) Update(ctx *gin.Context, entity any) (err error) {
	stmt := r.Database().Model(entity).Statement

	if entity, ok := entity.(types.IModelUpdatable); ok {
		stmt.Select(entity.Updatable(ctx))
	}
	if entity, ok := entity.(types.IModelUpdateClause); ok {
		AddClauses(stmt, entity.UpdateClause(ctx))
	}

	result := stmt.Updates(entity)
	if result.Error == nil && result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

// Delete removes an entity from the database.
func (r *DatabaseService) Delete(ctx *gin.Context, entity any) (err error) {
	stmt := r.Database().Model(entity).Statement
	if entity, ok := entity.(types.IModelDeleteClause); ok {
		AddClauses(stmt, entity.DeleteClause(ctx))
	}
	result := stmt.Delete(entity)
	if result.Error == nil && result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}
