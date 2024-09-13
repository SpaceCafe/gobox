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
		if len(fullPath) == 0 {
			return
		}
		if routerGroup != nil {
			fullPath = fullPath[len(routerGroup.BasePath()):]
		}
		ctx.Set("jwt/config", config)

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
				_ = ctx.Error(problems.ProblemJWTMissing)
			} else {
				config.Logger.Info(err)
				_ = ctx.Error(problems.ProblemJWTInvalid)
			}
			ctx.Abort()
			return
		}

		ctx.Set("jwt/token", token)
		ctx.Next()
	}
}
