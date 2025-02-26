package types

import (
	"github.com/gin-gonic/gin"
)

// IModel defines the basic methods that a model should implement to interact with its ID.
type IModel interface {
	// GetID returns the unique identifier of the model.
	GetID() string

	// SetID sets the unique identifier for the model.
	SetID(id string)
}

// IModelBeforeBind defines the method that should be called before binding data to the model.
type IModelBeforeBind interface {
	// BeforeBind is called before data is bound to the model from a request context.
	// This could be used to set default values.
	BeforeBind(ctx *gin.Context) (err error)
}

// IModelAfterBind defines the method that should be called after binding data to the model.
type IModelAfterBind interface {
	// AfterBind is called after data has been bound to the model from a request context.
	// This could be used to validate or modify the data before it is processed further,
	// but also to ensure that readonly fields are not processed.
	AfterBind(ctx *gin.Context) (err error)
}

// IModelBeforeRender defines the method that should be called before rendering the model.
type IModelBeforeRender interface {
	// BeforeRender is called before the model is rendered in a response context.
	// This could be used to modify or mask sensitive data before it is sent to the client.
	BeforeRender(ctx *gin.Context) (err error)
}

// IModelFilterable defines the method that should return filterable fields for the model.
type IModelFilterable interface {
	// Filterable returns a map of fields that can be used to filter queries on the model.
	// This ensures that only displayed values are filtered, preventing exposure of sensitive data.
	Filterable(ctx *gin.Context) map[string]struct{}
}

// IModelSortable defines the method that should return sortable fields for the model.
type IModelSortable interface {
	// Sortable returns a map of fields that can be used to sort queries on the model.
	// This ensures that only displayed values are sorted, preventing exposure of sensitive data.
	Sortable(ctx *gin.Context) map[string]struct{}
}

// IModelReadable defines the method that should return readable fields for the model.
// This ensures that only specified fields can be read, preventing loss of sensitive data.
type IModelReadable interface {
	// Readable returns a slice of field names that are allowed to be read in the model.
	Readable(ctx *gin.Context) []string
}

// IModelUpdatable defines the method that should return updatable fields for the model.
// This ensures that only specified fields can be updated, preventing modification of sensitive data.
type IModelUpdatable interface {
	// Updatable returns a slice of field names that are allowed to be updated in the model.
	Updatable(ctx *gin.Context) []string
}
