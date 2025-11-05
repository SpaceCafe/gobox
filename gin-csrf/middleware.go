package csrf

import (
	"net/http"

	"github.com/gin-gonic/gin"
	problems "github.com/spacecafe/gobox/gin-problems"
)

const (
	ContextKey = "urn:gobox:csrf:token"
)

// New creates a gin.HandlerFunc middleware to handle CSRF protection based on the provided configuration.
// This middleware implements the "double submit cookie" pattern. It validates the CSRF token from both
// header and cookie, aborts the request if validation fails.
// Only GET, HEAD, OPTIONS, and TRACE methods are allowed without token validation.
// We call `authentication.PrincipalFromContext` to retrieve the mandatory session ID from the context.
func New(cfg *Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		switch ctx.Request.Method {
		case http.MethodGet, http.MethodHead, http.MethodOptions, http.MethodTrace:
			ctx.Next()

			return
		}

		headerValue := ctx.Request.Header.Get(cfg.HeaderName)
		cookie, err := ctx.Request.Cookie(cfg.CookieName)

		if headerValue == "" || err != nil {
			problems.ProblemCSRFMissing.Abort(ctx)

			return
		}

		if headerValue != cookie.Value {
			problems.ProblemCSRFInvalid.Abort(ctx)

			return
		}

		err = ValidateToken(cfg, ctx, headerValue)
		if err != nil {
			problems.ProblemCSRFInvalid.Abort(ctx)

			return
		}

		ctx.Set(ContextKey, headerValue)
		ctx.Next()
	}
}

// TokenFromContext retrieves the CSRF token from the given Gin context if it exists.
// It returns the token and a boolean indicating whether the retrieval was successful.
func TokenFromContext(ctx *gin.Context) (string, bool) {
	value, ok := ctx.Get(ContextKey)
	if !ok {
		return "", false
	}

	token, ok := value.(string)

	return token, ok
}
