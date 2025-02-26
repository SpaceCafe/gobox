package types

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	// ResourceID is the constant for the resource identifier in the path parameter.
	ResourceID = "id"

	// ResourcePathWithoutID is the constant for the base path without an ID.
	ResourcePathWithoutID = "/"

	// ResourcePathWithID is the constant for the path with an ID parameter.
	ResourcePathWithID = ResourcePathWithoutID + ":" + ResourceID
)

// IResource defines the interface that all resources must implement,
// combining initialization and retrieval capabilities.
type IResource interface {
	IResourceInitializer
	IResourceGetter
}

// IResourceInitializer provides methods to set up a resource's dependencies and register it with the router.
type IResourceInitializer interface {
	// SetResources assigns the provided Resources object to the resource.
	// Additionally, it is responsible for initializing resource references in controllers, services, and views.
	SetResources(*Resources)

	// SetRouter sets the gin.IRouter for handling HTTP requests related to this resource.
	// The router is established by REST as a pre-configured group, incorporating the resource's designated path.
	SetRouter(router gin.IRouter)

	// Register registers the resource's routes and handlers with the router.
	// Must be called after SetRouter
	Register()
}

// IResourceGetter provides methods to retrieve information about a resource and its components.
type IResourceGetter interface {
	// Name returns the name of the resource.
	Name() string

	// Router returns the gin.IRouter associated with this resource.
	Router() gin.IRouter

	// Resource retrieves an IResourceGetter by name, allowing for nested resources.
	Resource() IResourceGetter

	// ResourceOf retrieves an IResourceGetter by name, allowing for nested resources.
	ResourceOf(name string) IResourceGetter

	// Controller returns the controller component of the resource.
	Controller() IController

	// Service returns the service component of the resource.
	Service() IService

	// View returns the view component of the resource.
	View() IView
}

// IResourceSetter provides a method to set an IResource.
type IResourceSetter interface {
	// SetResource assigns the provided IResource to this setter.
	SetResource(IResource)
}

// IResourceDatabase extends IResource with database access capabilities.
type IResourceDatabase interface {
	IResource
	// Database returns the gorm.DB instance associated with this resource for database operations.
	Database() *gorm.DB
}
