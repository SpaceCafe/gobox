package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spacecafe/gobox/gin-rest/types"
)

// DeleteFn is a generic struct that implements the types.ControllerDeleter interface.
type DeleteFn[T any] struct{}

// Delete is a method of DeleteFn that handles the deletion of a resource.
// It retrieves the resource ID from the request parameters, calls the service to delete the resource,
// and handles the response, setting the status to No Content if successful.
func (r *DeleteFn[T]) Delete(resource types.Resource[T]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param(types.PathParam)
		err := resource.GetService().(types.ServiceDeleter[T]).Delete(resource, id)
		if !HandleServiceError(ctx, err) {
			ctx.Status(http.StatusNoContent)
		}
	}
}
