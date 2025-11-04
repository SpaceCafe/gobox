package authentication

import (
	"github.com/gin-gonic/gin"
	problems "github.com/spacecafe/gobox/gin-problems"
)

var (
	_ Authenticator = (*TokenAuthenticator)(nil)
)

// TokenAuthenticator is responsible for authenticating requests using API tokens or basic authentication.
type TokenAuthenticator struct {
	cfg *Config
}

// NewTokenAuthenticator creates a new TokenAuthenticator with the given configuration.
func NewTokenAuthenticator(cfg *Config) *TokenAuthenticator {
	return &TokenAuthenticator{
		cfg: cfg,
	}
}

// Authenticate authenticates a request using the given context.
func (r *TokenAuthenticator) Authenticate(ctx *gin.Context) (Principal, error) {
	const prefix = "Token "
	auth := ctx.Request.Header.Get("Authorization")
	if !hasCaseInsensitivePrefix(auth, prefix) {
		return nil, ErrInvalidMethod
	}
	return r.cfg.Repository.GetByToken(auth[len(prefix):])
}

// Abort aborts the request with a 401 Unauthorized response and a WWW-Authenticate header.
func (r *TokenAuthenticator) Abort(ctx *gin.Context) {
	ctx.Header("WWW-Authenticate", "Token")
	problems.ProblemUnauthorized.Abort(ctx)
}
