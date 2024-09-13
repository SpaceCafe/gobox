package jsonapi

import (
	"github.com/gin-gonic/gin/render"
	"github.com/spacecafe/gobox/gin-rest/types"
)

// CRUD provides generic CRUD operations for a given entity type T.
// All operations return a render.Render to output the response as a JSON:API object.
type CRUD[T any] struct{}

// Create handles the creation of a new entity and returns a JSON:API response.
func (r *CRUD[T]) Create(resource types.Resource[T], entity *T, _ *types.ViewOptions) render.Render {
	return NewResponseFromEntity(resource, entity)
}

// Read handles reading an entity and returns a JSON:API response.

func (r *CRUD[T]) Read(resource types.Resource[T], entity *T, _ *types.ViewOptions) render.Render {
	return NewResponseFromEntity(resource, entity)
}

// List handles listing entities and returns a JSON:API response.
func (r *CRUD[T]) List(resource types.Resource[T], entities *[]T, options *types.ViewOptions) render.Render {
	return NewResponseFromEntities(resource, entities, options)
}

// Update handles updating an entity and returns a JSON:API response.
func (r *CRUD[T]) Update(resource types.Resource[T], _ string, entity *T, _ *types.ViewOptions) render.Render {
	return NewResponseFromEntity(resource, entity)
}

// Delete handles deleting an entity and returns a JSON:API response.
func (r *CRUD[T]) Delete(resource types.Resource[T], _ string, _ *types.ViewOptions) render.Render {
	return NewResponseFromEntity(resource, nil)
}

// Mime returns the MIME type for JSON:API responses.
func (r *CRUD[T]) Mime() string {
	return "application/vnd.api+json"
}
