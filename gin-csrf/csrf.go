package csrf

import (
	"errors"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
	problems "github.com/spacecafe/gobox/gin-problems"
	"github.com/spacecafe/gobox/httpserver"
)

func New(config *Config, serverConfig *httpserver.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		fullPath := ctx.FullPath()
		if len(fullPath) == 0 {
			return
		}
		ctx.Set("csrf/config", config)

		// Skip CSRF on excluded paths
		if slices.Contains(config.ExcludedRoutes, fullPath[len(serverConfig.BasePath):]) {
			ctx.Next()
			return
		}

		// Load and validate token from cookie
		token, err := NewTokenFromCookie(config, ctx)
		if (ctx.Request.Method == http.MethodGet || ctx.Request.Method == http.MethodHead || ctx.Request.Method == http.MethodOptions) && errors.Is(err, ErrCookieRetrieval) {
			if token, err = NewToken(config); err != nil {
				config.Logger.Warning(err)
				_ = ctx.Error(problems.ProblemCSRFMalfunction)
				ctx.Abort()
				return
			}
			http.SetCookie(ctx.Writer, &http.Cookie{
				Name:     config.CookieName,
				Value:    token.String(),
				Path:     config.Path,
				Domain:   config.Domain,
				Secure:   config.SecureCookie,
				HttpOnly: config.HTTPOnlyCookie,
				SameSite: config.sameSite,
			})
			ctx.Set("csrf/token", token)
			ctx.Next()
			return
		}

		if errors.Is(err, ErrCookieRetrieval) {
			_ = ctx.Error(problems.ProblemCSRFMissing)
			ctx.Abort()
			return
		}

		if err != nil {
			config.Logger.Info(err)
			_ = ctx.Error(problems.ProblemCSRFInvalid)
			ctx.Abort()
			return
		}

		// Validate CSRF token from header
		if !token.Compare(ctx.GetHeader(config.HeaderName)) {
			_ = ctx.Error(problems.ProblemCSRFInvalid)
			ctx.Abort()
			return
		}

		ctx.Set("csrf/token", token)
		ctx.Next()
	}
}
