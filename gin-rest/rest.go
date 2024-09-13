package rest

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/spacecafe/gobox/gin-rest/controller"
)

type REST struct {
	router      gin.IRouter
	schemaCache *sync.Map
}

func New(router gin.IRouter) *REST {

	// Check if router is nil.
	if router == nil {
		panic("router is nil")
	}

	// Initialize additional validators and check for errors.
	if err := controller.InitializeValidators(); err != nil {
		panic(err)
	}

	return &REST{
		router:      router,
		schemaCache: new(sync.Map),
	}
}
