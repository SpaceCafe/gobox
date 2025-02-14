package jwt

import (
	"errors"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5/request"
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

		// Skip JWT on excluded paths
		if slices.Contains(config.ExcludedRoutes, fullPath) {
			ctx.Next()
			return
		}

		// Load JWT token from header
		token, err := NewTokenFromHeader(config, ctx)
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
