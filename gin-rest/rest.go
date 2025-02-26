package rest

import (
	"github.com/gin-gonic/gin"
	authorization "github.com/spacecafe/gobox/gin-authorization"
	jwt "github.com/spacecafe/gobox/gin-jwt"
	"github.com/spacecafe/gobox/gin-rest/types"
)

// REST represents a RESTful API handler that manages resources and routes.
type REST struct {
	// router is a gin.IRouter used to define HTTP routes.
	router gin.IRouter

	// resources holds all registered resources for the REST API.
	resources *types.Resources
}

// New creates a new instance of REST with the provided router, JWT configuration, and authorization configuration.
func New(router *gin.RouterGroup, jwtConfig *jwt.Config, authorizationConfig *authorization.Config) (rest *REST) {
	// Check if router is nil.
	if router == nil {
		panic("router must be set")
	}

	// Initialize additional validators and check for errors.
	if err := InitializeValidators(); err != nil {
		panic(err)
	}

	rest = &REST{router: router, resources: &types.Resources{}}

	// Add JWT and authorization middlewares.
	rest.router.Use(jwt.New(jwtConfig, router))
	rest.router.Use(authorization.New(authorizationConfig, router))
	return
}

// Register adds a new resource to the REST API, setting up its routes and dependencies.
func (r *REST) Register(resource types.IResource) {
	if resource == nil {
		panic("resource must be set")
	}

	r.resources.Add(resource)
	resource.SetRouter(r.router.Group(resource.Name()))
	resource.SetResources(r.resources)
	resource.Controller().SetResource(resource)
	resource.Service().SetResource(resource)
	resource.View().SetResource(resource)
	resource.Router().Use(AcceptMiddleware(resource.View().SupportedMimeTypes()))
	resource.Register()
}
