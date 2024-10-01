package rest

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	authorization "github.com/spacecafe/gobox/gin-authorization"
	"github.com/spacecafe/gobox/gin-rest/controller"
	"github.com/spacecafe/gobox/gin-rest/service"
	"github.com/spacecafe/gobox/gin-rest/types"
	"github.com/spacecafe/gobox/gin-rest/view/json"
	"github.com/spacecafe/gobox/gin-rest/view/jsonapi"
	"github.com/spacecafe/gobox/gin-rest/view/yaml"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

const (
	// mimeSeparator is the string used to separate MIME types in the Allow and Accept headers.
	mimeSeparator = ", "

	ErrNoREST       = "rest must not be empty"
	ErrParsingModel = "could not parse model: %v"
	ErrNoRepository = "repository must not be empty"
	ErrNoService    = "service must not be empty if repository is not gorm.DB"
)

// Resource represents a generic resource with various components such as Controller, Service, Views, and Repository.
type Resource[T any] struct {

	// Controller handles the business logic for the resource.
	Controller any

	// Service provides additional functionalities for the resource.
	Service any

	// Views holds different views or representations of the resource depending on their mimetype.
	Views map[string]any

	// Repository manages the data persistence for the resource.
	Repository any

	// REST is a reference to the REST configuration.
	rest *REST

	// schema defines the structure of the resource.
	schema *schema.Schema

	// gormDB is the database connection using GORM.
	gormDB *gorm.DB

	// capabilities lists the capabilities of the resource.
	capabilities map[string][]string

	// basePath represents the absolute base URL of the resource.
	basePath *url.URL

	// namingStrategy is used to define naming conventions for the resource.
	namingStrategy *types.NamingStrategy

	// Group is a reference to the used gin.RouterGroup of the resource.
	Group *gin.RouterGroup

	// Authorization indicates whether the resource requires authorization.
	Authorization bool
}

// Apply initializes the Resource with the given REST configuration.
func (r *Resource[T]) Apply(rest *REST) {
	var err error

	// Check if REST object is nil.
	if rest == nil {
		panic(ErrNoREST)
	}
	r.rest = rest
	r.capabilities = make(map[string][]string)

	// Set default controller if nil.
	if r.Controller == nil {
		r.Controller = new(controller.CRUD[T])
	}

	// Set default view if nil.
	if r.Views == nil || len(r.Views) == 0 {
		viewJSON := json.CRUD[T]{}
		viewJSONAPI := jsonapi.CRUD[T]{}
		viewYAML := yaml.CRUD[T]{}
		r.Views = map[string]any{
			viewJSON.Mime():    &viewJSON,
			viewJSONAPI.Mime(): &viewJSONAPI,
			viewYAML.Mime():    &viewYAML,
		}
	}

	// Check if GetRepository is nil or of type gorm.DB and
	// set default service if nil and GetRepository is of type gorm.DB.
	if r.Repository == nil {
		panic(ErrNoRepository)
	}
	if db, ok := r.Repository.(*gorm.DB); ok {
		r.gormDB = db
		if r.Service == nil {
			r.Service = new(service.CRUD[T])
		}
	} else if r.Service == nil {
		panic(ErrNoService)
	}

	// Parse model schema.
	r.schema, err = r.parseModel()
	if err != nil {
		panic(fmt.Sprintf(ErrParsingModel, err))
	}

	r.applyRoutes()
}

// BasePath returns the absolute base path to the resource.
func (r *Resource[T]) BasePath() *url.URL {
	return r.basePath
}

// DB returns the gorm.DB connection used by the resource, if set.
func (r *Resource[T]) DB() *gorm.DB {
	return r.gormDB
}

// GetController returns the controller instance associated with the resource.
// This is often a handler function that handles HTTP requests or similar events.
func (r *Resource[T]) GetController() any {
	return r.Controller
}

// GetRepository returns the repository instance associated with the resource.
// This may handle data storage, querying, or other database operations.
func (r *Resource[T]) GetRepository() any {
	return r.Repository
}

// GetService returns the service instance associated with the resource.
// This may provide business logic or data processing functionality.
func (r *Resource[T]) GetService() any {
	return r.Service
}

// GetViews returns the view instances associated with the resource.
// They may be responsible for rendering UI templates or similar tasks.
func (r *Resource[T]) GetViews() map[string]any {
	return r.Views
}

// GetGroup returns the used gin.RouterGround of the resource.
func (r *Resource[T]) GetGroup() *gin.RouterGroup {
	return r.Group
}

// HasField checks if the schema has a field with the given name.
func (r *Resource[T]) HasField(name string) bool {
	_, ok := r.schema.FieldsByDBName[name]
	return ok
}

// HasReadableField checks if the schema has a field with the given name and read permission.
func (r *Resource[T]) HasReadableField(name string) bool {
	_, ok := r.schema.FieldsByDBName[name]
	return ok && r.schema.FieldsByDBName[name].Readable
}

// Name returns the name of the resource schema.
func (r *Resource[T]) Name() string {
	return r.schema.Table
}

// NamingStrategy returns the custom schema.Namer implementation.
func (r *Resource[T]) NamingStrategy() schema.Namer {
	return r.namingStrategy
}

// PrimaryField returns the database name of the primary field of the resource.
func (r *Resource[T]) PrimaryField() string {
	return r.schema.PrioritizedPrimaryField.DBName
}

// PrimaryValue returns the value of the primary field of the given entity.
func (r *Resource[T]) PrimaryValue(entity *T) string {
	value := reflect.ValueOf(*entity)
	field := value.FieldByName(r.schema.PrioritizedPrimaryField.Name)
	if !field.IsValid() {
		return ""
	}
	return field.String()
}

// Schema returns the resource schema.
func (r *Resource[T]) Schema() *schema.Schema {
	return r.schema
}

// applyRoutes registers the routes for the resource using the REST router.
func (r *Resource[T]) applyRoutes() {
	authFunc := noAuthorization
	if r.Authorization {
		authFunc = authorization.RequireAuthorization
	}

	// Set absolute base path of this resource.
	r.Group = r.rest.router.Group("/" + r.Name())
	r.basePath = &url.URL{Path: r.Group.BasePath()}

	if creator, ok := r.Controller.(types.ControllerCreator[T]); ok && checkServiceAndViews[T, types.ServiceCreator[T], types.ViewCreator[T]](r) {
		r.Group.POST(types.PathWithoutID, authFunc(r.Name(), authorization.CreateAction), creator.Create(r))
		r.capabilities[types.PathWithoutID] = append(r.capabilities[types.PathWithoutID], http.MethodPost)
	}

	if reader, ok := r.Controller.(types.ControllerReader[T]); ok && checkServiceAndViews[T, types.ServiceReader[T], types.ViewReader[T]](r) {
		r.Group.GET(types.PathWithID, authFunc(r.Name(), authorization.ReadAction), reader.Read(r))
		r.capabilities[types.PathWithID] = append(r.capabilities[types.PathWithID], http.MethodGet)
	}

	if lister, ok := r.Controller.(types.ControllerLister[T]); ok && checkServiceAndViews[T, types.ServiceLister[T], types.ViewLister[T]](r) {
		r.Group.GET(types.PathWithoutID, authFunc(r.Name(), authorization.ListAction), lister.List(r))
		r.capabilities[types.PathWithoutID] = append(r.capabilities[types.PathWithoutID], http.MethodGet)
	}

	if updater, ok := r.Controller.(types.ControllerUpdater[T]); ok && checkServiceAndViews[T, types.ServiceUpdater[T], types.ViewUpdater[T]](r) {
		r.Group.PUT(types.PathWithID, authFunc(r.Name(), authorization.UpdateAction), updater.Update(r, false))
		r.Group.PATCH(types.PathWithID, authFunc(r.Name(), authorization.UpdateAction), updater.Update(r, true))
		r.capabilities[types.PathWithID] = append(r.capabilities[types.PathWithID], http.MethodPut, http.MethodPatch)
	}

	if deleter, ok := r.Controller.(types.ControllerDeleter[T]); ok && checkServiceAndViews[T, types.ServiceDeleter[T], types.ViewDeleter[T]](r) {
		r.Group.DELETE(types.PathWithID, authFunc(r.Name(), authorization.DeleteAction), deleter.Delete(r))
		r.capabilities[types.PathWithID] = append(r.capabilities[types.PathWithID], http.MethodDelete)
	}

	if capabilitator, ok := r.Controller.(types.ControllerCapabilitator); ok {
		for relPath, allows := range r.capabilities {
			r.Group.OPTIONS(relPath, authFunc(r.Name(), authorization.CapabilitiesAction), capabilitator.Capability(strings.Join(allows, mimeSeparator), joinMapKeys(r.Views, mimeSeparator)))
		}
	}
}

// checkServiceAndViews checks if the resource's service and views implement the specified interfaces.
// T: The type of the resource.
// S: The interface that the service should implement.
// V: The interface that the views should implement.
func checkServiceAndViews[T any, S any, V any](resource *Resource[T]) bool {
	if _, ok := resource.Service.(S); ok {
		for _, view := range resource.Views {
			if _, ok = view.(V); !ok {
				return false
			}
		}
		return true
	}
	return false
}

// parseModel parses the model schema using the provided schema cache and own naming strategy.
func (r *Resource[T]) parseModel() (*schema.Schema, error) {
	var model T
	r.namingStrategy = &types.NamingStrategy{}
	return schema.Parse(model, r.rest.schemaCache, r.namingStrategy)
}

func noAuthorization(string, authorization.Action) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
	}
}
