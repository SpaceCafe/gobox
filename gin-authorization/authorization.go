package authorization

import (
	"slices"

	"github.com/gin-gonic/gin"
)

// New creates an authorization middleware handler function for the given configuration and router group.
// It checks if the request path is excluded from authorization. If not, it generates a new set of authorizations
// based on the configuration and context, stores them in the context, and proceeds to the next handler.
func New(config *Config, routerGroup *gin.RouterGroup) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		fullPath := ctx.FullPath()
		if fullPath == "" {
			return
		}
		if routerGroup != nil {
			fullPath = fullPath[len(routerGroup.BasePath()):]
		}

		// Skip authorization checks on excluded paths
		if slices.Contains(config.ExcludedRoutes, fullPath) {
			ctx.Next()
			return
		}

		authorizations := NewAuthorizations(config, ctx)
		SetAuthorizations(authorizations, ctx)
		ctx.Next()
	}
}
