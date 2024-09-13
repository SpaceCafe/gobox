package types

import (
	"github.com/gin-gonic/gin/render"
)

// ViewCreator is a generic interface for creating views.
// It defines a single method Create which takes a resource, an entity, and view options, and returns a render.Render.
type ViewCreator[T any] interface {
	Create(resource Resource[T], entity *T, options *ViewOptions) render.Render
}

// ViewDeleter is a generic interface for deleting views.
// It defines a single method Delete which takes a resource, an ID, and view options, and returns a render.Render.
type ViewDeleter[T any] interface {
	Delete(resource Resource[T], id string, options *ViewOptions) render.Render
}

// ViewLister is a generic interface for listing views.
// It defines a single method List which takes a resource, a slice of entities, and view options, and returns a render.Render.
type ViewLister[T any] interface {
	List(resource Resource[T], entities *[]T, options *ViewOptions) render.Render
}

// ViewReader is a generic interface for reading views.
// It defines a single method Read which takes a resource, an entity, and view options, and returns a render.Render.
type ViewReader[T any] interface {
	Read(resource Resource[T], entity *T, options *ViewOptions) render.Render
}

// ViewUpdater is a generic interface for updating views.
// It defines a single method Update which takes a resource, an ID, an entity, and view options, and returns a render.Render.
type ViewUpdater[T any] interface {
	Update(resource Resource[T], id string, entity *T, options *ViewOptions) render.Render
}
