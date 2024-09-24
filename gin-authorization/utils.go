package authorization

import (
	"github.com/gin-gonic/gin"
	problems "github.com/spacecafe/gobox/gin-problems"
)

// IsAuthorized checks if a given action on a resource is authorized for the current context.
func IsAuthorized(ctx *gin.Context, resource string, action Action) bool {
	if authorizations, ok := ctx.Get("authorization/authorizations"); ok {
		return authorizations.(Authorizations).IsAuthorized(resource, action)
	}
	return false
}

func RequireAuthorization(resource string, action Action) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if IsAuthorized(ctx, resource, action) {
			ctx.Next()
			return
		}
		_ = ctx.Error(problems.ProblemInsufficientPermission)
		ctx.Abort()
	}
}
