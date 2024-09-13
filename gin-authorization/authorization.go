package authorization

import (
	"slices"

	"github.com/gin-gonic/gin"
	problems "github.com/spacecafe/gobox/gin-problems"
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
		ctx.Set("authorization/authorizations", authorizations)

		for _, resourceAction := range NewResourceActions(ctx, fullPath) {
			if !authorizations.IsAuthorized(resourceAction.Resource, resourceAction.Action) {
				_ = ctx.Error(problems.ProblemInsufficientPermission)
				ctx.Abort()
				return
			}
		}

		ctx.Next()
	}
}
