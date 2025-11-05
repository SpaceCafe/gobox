package authentication

import (
	"github.com/gin-gonic/gin"
	problems "github.com/spacecafe/gobox/gin-problems"
)

var _ Authenticator = (*BasicAuthenticator)(nil)

// BasicAuthenticator is responsible for authenticating requests using basic authentication.
type BasicAuthenticator struct {
	cfg *Config
}

// NewBasicAuthenticator creates a new BasicAuthenticator with the given configuration.
func NewBasicAuthenticator(cfg *Config) *BasicAuthenticator {
	return &BasicAuthenticator{
		cfg: cfg,
	}
}

// Abort aborts the request with a 401 Unauthorized response and a WWW-Authenticate header.
func (r *BasicAuthenticator) Abort(ctx *gin.Context) {
	ctx.Header("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
	problems.ProblemUnauthorized.Abort(ctx)
}

// Authenticate authenticates a request using the given context.
//
//nolint:ireturn // Principal is implemented by the repository.
func (r *BasicAuthenticator) Authenticate(ctx *gin.Context) (Principal, error) {
	id, password, ok := ctx.Request.BasicAuth()
	if !ok {
		return nil, ErrInvalidMethod
	}

	//nolint:wrapcheck // wrap check is not relevant here.
	return r.cfg.Repository.GetByCredentials(id, password)
}
