package types

import (
	"github.com/gin-gonic/gin"
)

// IView defines the contract for view operations in a RESTful API,
// responsible for rendering responses.
type IView interface {
	// IResourceGetter is embedded to include methods related to retrieving resources.
	IResourceGetter

	// IResourceSetter is embedded to include methods related to modifying resources.
	IResourceSetter

	// Create handles the creation of a new resource and returns a IRender object.
	Create(ctx *gin.Context, entity any) IRender

	// Read retrieves an existing resource and returns a IRender object.
	Read(ctx *gin.Context, entity any) IRender

	// List retrieves a collection of resources and returns a IRender object.
	List(ctx *gin.Context, entities any) IRender

	// Update retrieves an existing resource and returns a IRender object.
	Update(ctx *gin.Context, entity any) IRender

	// Delete retrieves an existing resource and returns a IRender object.
	Delete(ctx *gin.Context, entity any) IRender

	// SupportedMimeTypes returns a pointer to a map containing supported MIME types.
	SupportedMimeTypes() map[string]any
}
