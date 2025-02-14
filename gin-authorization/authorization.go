package authorization

import (
	"slices"

	"github.com/gin-gonic/gin"
)

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
		SetAuthorizations(&authorizations, ctx)
		ctx.Next()
	}
}
