package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	problems "github.com/spacecafe/gobox/gin-problems"
	"github.com/spacecafe/gobox/gin-rest/types"
)

// UpdateFn is a generic struct that implements the types.ControllerUpdater interface.
type UpdateFn[T any] struct{}

// Update is a method of UpdateFn that handles the updating of a resource.
// It retrieves the resource ID from the request parameters, binds the JSON payload to an entity,
// calls the service to update the entity, and handles the response, rendering the view if successful.
func (r *UpdateFn[T]) Update(resource types.Resource[T], partially bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var entity T
		id := ctx.Param(types.PathParam)
		if err := ctx.ShouldBindJSON(&entity); err != nil {
			HandleError(ctx, err, problems.ProblemBadRequest)
			return
		}

		_, err := resource.GetService().(types.ServiceUpdater[T]).Update(resource, partially, id, &entity)
		if !HandleServiceError(ctx, err) {
			ctx.Status(http.StatusNoContent)
		}
	}
}
