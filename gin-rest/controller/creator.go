package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spacecafe/gobox/gin-rest/types"
)

// CreateFn is a generic struct that implements the types.ControllerCreator interface.
type CreateFn[T any] struct{}

// Create is a method of CreateFn that handles the creation of a resource.
// It binds the JSON payload to an entity, calls the service to create the entity,
// and handles the response, including setting the Location header and rendering the view.
func (r *CreateFn[T]) Create(resource types.Resource[T]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var entity T
		if HandleControllerError(ctx, ctx.ShouldBindJSON(&entity)) {
			return
		}

		createdEntity, err := resource.GetService().(types.ServiceCreator[T]).Create(resource, NewServiceOptions(ctx), &entity)
		if view := GetView[T](ctx, resource); !HandleServiceError(ctx, err) && view != nil {
			url := ctx.Request.URL.JoinPath(resource.PrimaryValue(createdEntity))
			ctx.Header("Location", url.String())
			ctx.Render(http.StatusCreated, view.(types.ViewCreator[T]).Create(resource, createdEntity, &types.ViewOptions{ServiceMeta: types.ServiceMeta{Total: 1}}))
		}
	}
}
