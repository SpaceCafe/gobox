package jwt

import (
	"errors"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5/request"
	problems "github.com/spacecafe/gobox/gin-problems"
)

// New creates a JWT middleware handler function for the given configuration and router group.
// It checks if the request path is excluded from JWT validation. If not, it attempts to load
// a JWT token from the request. If successful, the token is stored in the context; otherwise,
// appropriate error responses are sent.
func New(config *Config, routerGroup *gin.RouterGroup) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		fullPath := ctx.FullPath()
		if fullPath == "" {
			return
		}
		if routerGroup != nil {
			fullPath = fullPath[len(routerGroup.BasePath()):]
		}

		// Skip JWT on excluded paths
		if slices.Contains(config.ExcludedRoutes, fullPath) {
			ctx.Next()
			return
		}

		// Load JWT token from request
		token, err := NewTokenFromRequest(config, ctx)
		if err != nil {
			ctx.Header("WWW-Authenticate", "Bearer realm=\"JWT\", charset=\"UTF-8\"")
			if errors.Is(err, request.ErrNoTokenInRequest) {
				problems.ProblemJWTMissing.Abort(ctx)
			} else {
				_ = ctx.Error(err)
				problems.ProblemJWTInvalid.Abort(ctx)
			}
			return
		}

		SetToken(token, ctx)
		ctx.Next()
	}
}
