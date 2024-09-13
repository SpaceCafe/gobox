package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spacecafe/gobox/gin-rest/types"
)

// ListFn is a generic struct that implements the types.ControllerLister interface.
type ListFn[T any] struct{}

// List is a method of ListFn that handles the listing of resources.
// It calls the service to list the entities, retrieves metadata, and handles the response,
// rendering the view if successful.
func (r *ListFn[T]) List(resource types.Resource[T]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		entities, meta, err := resource.GetService().(types.ServiceLister[T]).List(resource, NewRequestParams(ctx, resource.HasReadableField).Options())
		if view := getView[T](ctx, resource); !handleServiceError(ctx, err) && view != nil {
			ctx.Render(http.StatusOK, view.(types.ViewLister[T]).List(resource, entities, &types.ViewOptions{ServiceMeta: *meta}))
		}
	}
}
