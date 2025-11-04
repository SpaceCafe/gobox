package authentication

import (
	"errors"

	"github.com/gin-gonic/gin"
)

const (
	PrincipalContextKey = "urn:gobox:authentication:principal"
)

// New returns a gin.HandlerFunc that authenticates requests based on the provided configuration.
func New(cfg *Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var i int
		for i = range cfg.Authenticators {
			principal, err := cfg.Authenticators[i].Authenticate(ctx)
			if err == nil && principal != nil {
				ctx.Set(PrincipalContextKey, principal)
				ctx.Next()
				return
			}

			if errors.Is(err, ErrInvalidMethod) {
				continue
			}

			_ = ctx.Error(err)
			break
		}

		cfg.Authenticators[i].Abort(ctx)
	}
}

// PrincipalFromContext retrieves the Principal from the given Gin context if it exists.
// It returns the Principal and a boolean indicating whether the retrieval was successful.
func PrincipalFromContext(ctx *gin.Context) (Principal, bool) {
	value, ok := ctx.Get(PrincipalContextKey)
	if !ok {
		return nil, false
	}
	principal, ok := value.(Principal)
	return principal, ok
}
