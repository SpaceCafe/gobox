package authentication

import (
	"github.com/gin-gonic/gin"
)

// Authenticator defines an interface for authenticating users or services in the context of a Gin web application.
type Authenticator interface {
	// Authenticate authenticates the user in the context of a Gin web application.
	// It returns a Principal compliant struct on success or an error if authentication fails.
	// The Principal is stored by this middleware into the gin.Context later.
	Authenticate(ctx *gin.Context) (Principal, error)

	// Abort is called when authentication fails to provide a custom response.
	Abort(ctx *gin.Context)
}
