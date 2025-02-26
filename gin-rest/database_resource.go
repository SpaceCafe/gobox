package rest

import (
	"github.com/spacecafe/gobox/gin-rest/types"
	"gorm.io/gorm"
)

// Ensure DatabaseResource implements the IResourceDatabase interface.
var _ types.IResourceDatabase = (*DatabaseResource)(nil)

// DatabaseResource represents a resource that interacts with a database using GORM.
type DatabaseResource struct {
	// BaseResource is an embedded field representing the base resource functionalities.
	BaseResource

	// database holds the reference to the GORM DB instance used for database operations.
	database *gorm.DB
}

// NewDatabaseResource creates a new instance of DatabaseResource.
// It requires a controller, service, view, and a non-nil database connection.
func NewDatabaseResource[T any](controller types.IController, service types.IService, view types.IView, database *gorm.DB) *DatabaseResource {
	if database == nil {
		panic("database must be set")
	}
	if service == nil {
		service = &DatabaseService{}
	}
	res := &DatabaseResource{BaseResource: *NewResource[T](controller, service, view), database: database}

	var entity T
	stmt := &gorm.Statement{DB: database}
	if err := stmt.Parse(entity); err != nil {
		panic(err)
	}
	res.name = stmt.Schema.Table
	return res
}

// Database returns the GORM DB instance associated with this resource.
func (r *DatabaseResource) Database() *gorm.DB {
	return r.database
}
