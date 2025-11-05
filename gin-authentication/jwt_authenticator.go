package authentication

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/spacecafe/gobox/gin-authentication/jwt"
	problems "github.com/spacecafe/gobox/gin-problems"
)

var (
	_ Authenticator = (*JWTAuthenticator)(nil)
	_ Principal     = (*jwt.Claims)(nil)
)

// JWTAuthenticator is responsible for authenticating requests using bearer authentication.
type JWTAuthenticator struct {
	cfg *Config
}

// NewJWTAuthenticator creates a new JWTAuthenticator with the given configuration.
func NewJWTAuthenticator(cfg *Config) *JWTAuthenticator {
	return &JWTAuthenticator{
		cfg: cfg,
	}
}

// Abort aborts the request with a 401 Unauthorized response and a WWW-Authenticate header.
func (r *JWTAuthenticator) Abort(ctx *gin.Context) {
	ctx.Header("WWW-Authenticate", "JWT")

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
func (r *JWTAuthenticator) Authenticate(ctx *gin.Context) (Principal, error) {
	cookie, err := ctx.Request.Cookie(r.cfg.JWT.CookieName)
	if err != nil {
		return nil, ErrInvalidMethod
	}

	token, err := jwt.NewFromString(r.cfg.JWT, cookie.Value, jwt.AccessToken)
	if err != nil {
		return nil, problems.ProblemJWTInvalid.WithError(err)
	}

	return token.Claims(), nil
}
