package authentication

import (
	"github.com/gin-gonic/gin"
	problems "github.com/spacecafe/gobox/gin-problems"
)

var _ Authenticator = (*TokenAuthenticator)(nil)

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

// Abort aborts the request with a 401 Unauthorized response and a WWW-Authenticate header.
func (r *TokenAuthenticator) Abort(ctx *gin.Context) {
	ctx.Header("WWW-Authenticate", "Token")
	problems.ProblemUnauthorized.Abort(ctx)
}

// Authenticate authenticates a request using the given context.
//
//nolint:ireturn // Principal is implemented by the repository.
func (r *TokenAuthenticator) Authenticate(ctx *gin.Context) (Principal, error) {
	const prefix = "Token "

	auth := ctx.Request.Header.Get("Authorization")
	if !hasCaseInsensitivePrefix(auth, prefix) {
		return nil, ErrInvalidMethod
	}

	//nolint:wrapcheck // wrap check is not relevant here.
	return r.cfg.Repository.GetByToken(auth[len(prefix):])
}
