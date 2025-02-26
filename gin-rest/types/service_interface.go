package types

import (
	"github.com/gin-gonic/gin"
)

// IService defines the contract for service operations in a RESTful API.
type IService interface {
	// IResourceGetter is embedded to include methods related to retrieving resources.
	IResourceGetter

	// IResourceSetter is embedded to include methods related to modifying resources.
	IResourceSetter

	// Create handles the creation of a new resource.
	Create(ctx *gin.Context, entity any) (err error)

	// Read retrieves an existing resource.
	Read(ctx *gin.Context, entity any) (err error)

	// List retrieves a collection of resources.
	List(ctx *gin.Context, entities any) (err error)

	// Update modifies an existing resource.
	Update(ctx *gin.Context, entity any) (err error)

	// Delete removes a resource.
	Delete(ctx *gin.Context, entity any) (err error)
}
