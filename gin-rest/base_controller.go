package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spacecafe/gobox/gin-rest/types"
)

// Ensure BaseController implements types.IController interface.
var _ types.IController = (*BaseController[types.IModel])(nil)

// BaseController is a generic controller for handling CRUD operations on resources of type T.
type BaseController[T any] struct {
	types.IResourceGetter
}

// SetResource sets the resource for this controller.
func (r *BaseController[T]) SetResource(resource types.IResource) {
	r.IResourceGetter = resource
}

// Create returns a gin.HandlerFunc that handles HTTP POST requests to create a new entity of type T.
func (r *BaseController[T]) Create() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var entity T
		if BindJSON(ctx, &entity) || HandleError(ctx, r.Service().Create(ctx, &entity)) || BeforeRenderHook(ctx, &entity) {
			return
		}
		if entity, ok := any(entity).(types.IModel); ok {
			ctx.Header("Location", ctx.Request.URL.JoinPath(entity.GetID()).String())
		}
		ctx.Render(http.StatusCreated, r.View().Create(ctx, &entity))
	}
}

// Read returns a gin.HandlerFunc that handles HTTP GET requests to read an entity of type T.
func (r *BaseController[T]) Read() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var entity T
		if BindID(ctx, &entity) || AfterBindHook(ctx, &entity) || HandleError(ctx, r.Service().Read(ctx, &entity)) || BeforeRenderHook(ctx, &entity) {
			return
		}
		ctx.Render(http.StatusOK, r.View().Read(ctx, &entity))
	}
}

// List returns a gin.HandlerFunc that handles HTTP GET requests to list entities of type T.
func (r *BaseController[T]) List() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var entity T
		var entities []T
		if HandleError(ctx, ParseListOptions(ctx, &entity)) || HandleError(ctx, r.Service().List(ctx, &entities)) {
			return
		}
		for _, entity := range entities {
			if BeforeRenderHook(ctx, entity) {
				return
			}
		}
		ctx.Render(http.StatusOK, r.View().List(ctx, &entities))
	}
}

// Update returns a gin.HandlerFunc that handles HTTP PUT requests to update an entity of type T.
func (r *BaseController[T]) Update() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var entity T
		if BindID(ctx, &entity) || BindJSON(ctx, &entity) || HandleError(ctx, r.Service().Update(ctx, &entity)) || BeforeRenderHook(ctx, &entity) {
			return
		}
		ctx.Status(http.StatusNoContent)
	}
}

// Delete returns a gin.HandlerFunc that handles HTTP DELETE requests to delete an entity of type T.
func (r *BaseController[T]) Delete() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var entity T
		if BindID(ctx, &entity) || AfterBindHook(ctx, &entity) || HandleError(ctx, r.Service().Delete(ctx, &entity)) {
			return
		}
		ctx.Status(http.StatusNoContent)
	}
}
