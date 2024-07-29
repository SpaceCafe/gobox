package csrf

import (
	"github.com/gin-gonic/gin"
)

// HeaderToken retrieves the CSRF token and its associated header name from the context.
// It returns the header name and the encoded token if they are present in the context,
// otherwise it returns empty strings.
func HeaderToken(ctx *gin.Context) (headerName, encodedToken string) {
	config, okConfig := ctx.Get("csrf/config")
	token, okToken := ctx.Get("csrf/token")

	if okConfig && okToken {
		return config.(*Config).HeaderName, token.(*Token).encodedToken
	}
	return "", ""
}
