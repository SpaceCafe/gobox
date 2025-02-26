package types

import (
	"github.com/gin-gonic/gin"
)

// IController defines the contract for a controller in a RESTful API.
// It includes methods for handling CRUD operations and embedding interfaces for resource getting and setting.
type IController interface {
	// IResourceGetter is embedded to include methods related to retrieving resources.
	IResourceGetter

	// IResourceSetter is embedded to include methods related to modifying resources.
	IResourceSetter

	// Create returns a gin.HandlerFunc that handles the creation of a new resource.
	Create() gin.HandlerFunc

	// Read returns a gin.HandlerFunc that handles reading a specific resource.
	Read() gin.HandlerFunc

	// List returns a gin.HandlerFunc that handles listing multiple resources.
	List() gin.HandlerFunc

	// Update returns a gin.HandlerFunc that handles updating an existing resource.
	Update() gin.HandlerFunc

	// Delete returns a gin.HandlerFunc that handles deleting a resource.
	Delete() gin.HandlerFunc
}
