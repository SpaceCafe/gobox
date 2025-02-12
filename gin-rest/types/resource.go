package types

import (
	"net/url"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

const (
	// PathParam is the constant for the resource identifier in the path parameter.
	PathParam = "id"

	// PathWithoutID is the constant for the base path without an ID.
	PathWithoutID = "/"

	// PathWithID is the constant for the path with an ID parameter.
	PathWithID = PathWithoutID + ":" + PathParam
)

// Resource is a generic interface for a resource with type T.
type Resource[T any] interface {

	// BasePath returns the absolute base path to the resource.
	BasePath() *url.URL

	// DB returns the gorm database instance.
	DB() *gorm.DB

	// GetController returns the controller associated with the resource.
	GetController() any

	// GetRepository returns the repository associated with the resource.
	GetRepository() any

	// GetService returns the service associated with the resource.
	GetService() any

	// GetViews returns a map of views associated with the resource.
	GetViews() map[string]any

	// GetGroup returns the used gin.RouterGroup of the resource.
	GetGroup() *gin.RouterGroup

	// HasField checks if the resource has a specific field.
	HasField(name string) bool

	// HasReadableField checks if the resource has a specific readable field.
	HasReadableField(name string) bool

	// Name returns the name of the resource.
	Name() string

	// NamingStrategy returns a schema.Namer implementation.
	NamingStrategy() schema.Namer

	// PrimaryField returns the primary field of the resource.
	PrimaryField() string

	// PrimaryValue returns the primary value of the entity.
	PrimaryValue(entity *T) string

	// Schema returns the resource schema.
	Schema() *schema.Schema
}
