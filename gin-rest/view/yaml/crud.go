package yaml

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/gin-gonic/gin/render"
	"github.com/spacecafe/gobox/gin-rest/types"
)

// CRUD provides generic CRUD operations for a given entity type T.
// All operations return a render.Render to output the response as a YAML object.
type CRUD[T any] struct{}

// Create handles the creation of a new entity and returns a YAML response.
func (r *CRUD[T]) Create(_ types.Resource[T], entity *T, options *types.ViewOptions) render.Render {
	return r.createResponse(&[]T{*entity}, options)
}

// Read handles reading an entity and returns a YAML response.
func (r *CRUD[T]) Read(_ types.Resource[T], entity *T, options *types.ViewOptions) render.Render {
	return r.createResponse(&[]T{*entity}, options)
}

// List handles listing entities and returns a YAML response.
func (r *CRUD[T]) List(_ types.Resource[T], entities *[]T, options *types.ViewOptions) render.Render {
	return r.createResponse(entities, options)
}

// Update handles updating an entity and returns a YAML response.
func (r *CRUD[T]) Update(_ types.Resource[T], _ string, entity *T, options *types.ViewOptions) render.Render {
	return r.createResponse(&[]T{*entity}, options)
}

// Delete handles deleting an entity and returns a YAML response.
func (r *CRUD[T]) Delete(_ types.Resource[T], _ string, options *types.ViewOptions) render.Render {
	return r.createResponse(&[]T{}, options)
}

// Mime returns the MIME type for YAML responses.
func (r *CRUD[T]) Mime() string {
	return binding.MIMEYAML
}

// createResponse creates a YAML response with pagination and entity data.
func (r *CRUD[T]) createResponse(entities *[]T, options *types.ViewOptions) render.Render {
	return &Response[T]{
		Page:       options.GetPage(),
		PageSize:   options.GetPageSize(),
		Total:      options.Total,
		TotalPages: options.GetTotalPages(),
		Data:       entities,
	}
}
