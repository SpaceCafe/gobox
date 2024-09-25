package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spacecafe/gobox/gin-rest/types"
)

// ReadFn is a generic struct that implements the types.ControllerReader interface.
type ReadFn[T any] struct{}

// Read is a method of ReadFn that handles the reading of a resource.
// It retrieves the resource ID from the request parameters, calls the service to read the resource,
// and handles the response, rendering the view if successful.
func (r *ReadFn[T]) Read(resource types.Resource[T]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param(types.PathParam)
		entity, err := resource.GetService().(types.ServiceReader[T]).Read(resource, id)
		if view := GetView[T](ctx, resource); !HandleServiceError(ctx, err) && view != nil {
			ctx.Render(http.StatusOK, view.(types.ViewReader[T]).Read(resource, entity, &types.ViewOptions{ServiceMeta: types.ServiceMeta{Total: 1}}))
		}
	}
}
