package rest

import (
	"sync"

	"github.com/gin-gonic/gin"
	authorization "github.com/spacecafe/gobox/gin-authorization"
	jwt "github.com/spacecafe/gobox/gin-jwt"
	"github.com/spacecafe/gobox/gin-rest/controller"
)

type REST struct {
	router      *gin.RouterGroup
	schemaCache *sync.Map
}

func New(router *gin.RouterGroup, jwtConfig *jwt.Config, authorizationConfig *authorization.Config) *REST {

	// Check if router is nil.
	if router == nil {
		panic("router is nil")
	}

	// Initialize additional validators and check for errors.
	if err := controller.InitializeValidators(); err != nil {
		panic(err)
	}

	// Add JWT and authorization middlewares.
	router.Use(jwt.New(jwtConfig, router))
	router.Use(authorization.New(authorizationConfig, router))

	return &REST{
		router:      router,
		schemaCache: new(sync.Map),
	}
}
