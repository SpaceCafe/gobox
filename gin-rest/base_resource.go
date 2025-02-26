package rest

import (
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	authorization "github.com/spacecafe/gobox/gin-authorization"
	"github.com/spacecafe/gobox/gin-rest/types"
)

// Ensure BaseResource implements the IResource interface.
var _ types.IResource = (*BaseResource)(nil)

// BaseResource represents a RESTful resource with associated controller, service, and view.
type BaseResource struct {
	// router is the Gin router used for registering routes.
	router gin.IRouter

	// controller handles business logic operations.
	controller types.IController

	// service provides data access and manipulation.
	service types.IService

	// view manages response formatting.
	view types.IView

	// name is the snake_case representation of the resource's model type.
	name string

	// resources holds a reference to the Resources manager.
	resources *types.Resources
}

// NewResource creates a new BaseResource instance with the provided controller, service, and view.
func NewResource[T any](controller types.IController, service types.IService, view types.IView) *BaseResource {
	var res BaseResource

	// Check if the model is a struct and set its name. Otherwise, panic.
	if reflect.TypeFor[T]().Kind() != reflect.Struct {
		panic("model must be a struct")
	}
	if a, b, ok := strings.Cut(reflect.TypeFor[T]().Name(), "."); ok {
		res.name = CamelToSnake(b)
	} else {
		res.name = CamelToSnake(a)
	}

	// Initialize the controller if it's not set.
	if controller == nil {
		res.controller = &Controller[T]{}
	} else {
		res.controller = controller
	}

	// Check if the service is set and panic if not.
	if service == nil {
		panic("service must be set")
	}
	res.service = service

	// Initialize the view if it's not set.
	if view == nil {
		res.view = &View{}
	} else {
		res.view = view
	}
	return &res
}

// SetResources assigns a Resources manager to the BaseResource.
func (r *BaseResource) SetResources(resources *types.Resources) {
	if resources == nil {
		panic("resources is nil")
	}
	r.resources = resources
}

// SetRouter assigns a Gin router to the BaseResource for route registration.
func (r *BaseResource) SetRouter(router gin.IRouter) {
	if router == nil {
		panic("router is nil")
	}
	r.router = router
}

// Register sets up HTTP routes for the resource with appropriate authorization and controller methods.
func (r *BaseResource) Register() {
	r.Router().POST(types.ResourcePathWithoutID, authorization.RequireAuthorization(r.Name(), authorization.CreateAction), r.Controller().Create())
	r.Router().GET(types.ResourcePathWithID, authorization.RequireAuthorization(r.Name(), authorization.ReadAction), r.Controller().Read())
	r.Router().GET(types.ResourcePathWithoutID, authorization.RequireAuthorization(r.Name(), authorization.ListAction), r.Controller().List())
	r.Router().PUT(types.ResourcePathWithID, authorization.RequireAuthorization(r.Name(), authorization.UpdateAction), r.Controller().Update())
	r.Router().DELETE(types.ResourcePathWithID, authorization.RequireAuthorization(r.Name(), authorization.DeleteAction), r.Controller().Delete())
}

// Name returns the snake_case name of the resource.
func (r *BaseResource) Name() string {
	return r.name
}

// Router returns the Gin router associated with the resource.
func (r *BaseResource) Router() gin.IRouter {
	return r.router
}

// Resource retrieves a specific resource by name from the Resources manager.
func (r *BaseResource) Resource(name string) types.IResourceGetter {
	return r.resources.Get(name)
}

// Controller returns the controller associated with the resource.
func (r *BaseResource) Controller() types.IController {
	return r.controller
}

// Service returns the service associated with the resource.
func (r *BaseResource) Service() types.IService {
	return r.service
}

// View returns the view associated with the resource.
func (r *BaseResource) View() types.IView {
	return r.view
}
