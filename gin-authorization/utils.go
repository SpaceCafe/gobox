package authorization

import (
	"github.com/gin-gonic/gin"
	problems "github.com/spacecafe/gobox/gin-problems"
)

// GetAuthorizations retrieves the Authorizations object from the context.
// It returns an empty Authorizations object if not found.
func GetAuthorizations(ctx *gin.Context) Authorizations {
	if authorizations, ok := ctx.Get("authorization/authorizations"); ok {
		return authorizations.(Authorizations)
	}
	return Authorizations{}
}

// IsAuthorized checks if a given action on a resource is authorized for the current context.
func IsAuthorized(ctx *gin.Context, resource string, action Action) bool {
	return GetAuthorizations(ctx).IsAuthorized(resource, action)
}

// RequireAuthorization creates a middleware handler that ensures the request is authorized
// for the specified resource and action. If not authorized, it aborts the request with an error.
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
