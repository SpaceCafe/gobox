package types

import (
	"github.com/gin-gonic/gin"
)

// ControllerCapabilitator is an interface for setting capability headers.
// It defines a single method Capability which takes allow and accept strings and returns a gin.HandlerFunc.
type ControllerCapabilitator interface {
	Capability(allow, accept string) gin.HandlerFunc
}

// ControllerCreator is a generic interface for creating resources.
// It defines a single method Create which takes a resource and returns a gin.HandlerFunc.
type ControllerCreator[T any] interface {
	Create(resource Resource[T]) gin.HandlerFunc
}

// ControllerDeleter is a generic interface for deleting resources.
// It defines a single method Delete which takes a resource and returns a gin.HandlerFunc.
type ControllerDeleter[T any] interface {
	Delete(resource Resource[T]) gin.HandlerFunc
}

// ControllerLister is a generic interface for listing resources.
// It defines a single method List which takes a resource and returns a gin.HandlerFunc.
type ControllerLister[T any] interface {
	List(resource Resource[T]) gin.HandlerFunc
}

// ControllerReader is a generic interface for reading resources.
// It defines a single method Read which takes a resource and returns a gin.HandlerFunc.
type ControllerReader[T any] interface {
	Read(resource Resource[T]) gin.HandlerFunc
}

// ControllerUpdater is a generic interface for updating resources.
// It defines a single method Update which takes a resource and a boolean indicating partial update, and returns a gin.HandlerFunc.
type ControllerUpdater[T any] interface {
	Update(resource Resource[T], partially bool) gin.HandlerFunc
}
