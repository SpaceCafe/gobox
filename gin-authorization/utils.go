package authorization

import (
	"github.com/gin-gonic/gin"
	problems "github.com/spacecafe/gobox/gin-problems"
)

const (
	ContextKeyName = "authorization.authorizations"
)

// SetAuthorizations stores the provided Authorizations object in the context.
func SetAuthorizations(authorizations *Authorizations, ctx *gin.Context) {
	ctx.Set(ContextKeyName, authorizations)
}

// GetAuthorizations retrieves the Authorizations object from the context. If no object is found,
// it returns a new empty Authorizations object instance.
func GetAuthorizations(ctx *gin.Context) (authorizations *Authorizations) {
	if authorizationsRaw, ok := ctx.Get(ContextKeyName); ok {
		if authorizations, ok = authorizationsRaw.(*Authorizations); ok {
			return
		}
	}
	return &Authorizations{}
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
		problems.ProblemInsufficientPermission.Abort(ctx)
	}
}
