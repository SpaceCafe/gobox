package authorization

import (
	"github.com/gin-gonic/gin"
)

// IsAuthorized checks if a given action on a resource is authorized for the current context.
func IsAuthorized(ctx *gin.Context, resource string, action Action) bool {
	if authorizations, ok := ctx.Get("authorization/authorizations"); ok {
		return authorizations.(Authorizations).IsAuthorized(resource, action)
	}
	return false
}
