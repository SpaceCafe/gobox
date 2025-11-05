package authentication

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/spacecafe/gobox/gin-authentication/jwt"
	problems "github.com/spacecafe/gobox/gin-problems"
)

var (
	_ Authenticator = (*BearerAuthenticator)(nil)
	_ Principal     = (*jwt.Claims)(nil)
)

// BearerAuthenticator is responsible for authenticating requests using bearer authentication.
type BearerAuthenticator struct {
	cfg *Config
}

// NewBearerAuthenticator creates a new BearerAuthenticator with the given configuration.
func NewBearerAuthenticator(cfg *Config) *BearerAuthenticator {
	return &BearerAuthenticator{
		cfg: cfg,
	}
}

// Abort aborts the request with a 401 Unauthorized response and a WWW-Authenticate header.
func (r *BearerAuthenticator) Abort(ctx *gin.Context) {
	ctx.Header("WWW-Authenticate", "Bearer")

	err := ctx.Errors.Last()
	if err != nil {
		var p *problems.Problem
		if errors.As(err.Err, &p) {
			ctx.Abort()
		}
	}

	problems.ProblemJWTMissing.Abort(ctx)
}

// Authenticate authenticates a request using the given context.
//
//nolint:ireturn // Principal is implemented by the repository.
func (r *BearerAuthenticator) Authenticate(ctx *gin.Context) (Principal, error) {
	const prefix = "Bearer "

	auth := ctx.Request.Header.Get("Authorization")
	if !hasCaseInsensitivePrefix(auth, prefix) {
		return nil, ErrInvalidMethod
	}

	token, err := jwt.NewFromString(r.cfg.JWT, auth[len(prefix):], jwt.AccessToken)
	if err != nil {
		return nil, problems.ProblemJWTInvalid.WithError(err)
	}

	return token.Claims(), nil
}
