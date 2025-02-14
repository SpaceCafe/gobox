package jwt

import (
	"github.com/gin-gonic/gin"
)

const (
	ContextKeyName = "jwt.token"
)

// SetToken stores the provided JWT token in the context.
func SetToken(token *Token, ctx *gin.Context) {
	ctx.Set(ContextKeyName, token)
}

// GetToken retrieves the JWT token from the context. If no token is found,
// it returns a new empty Token instance.
func GetToken(ctx *gin.Context) (token *Token) {
	if tokenRaw, ok := ctx.Get(ContextKeyName); ok {
		if token, ok = tokenRaw.(*Token); ok {
			return
		}
	}
	return &Token{}
}

// GetClaims extracts the claims from the JWT token stored in the context.
func GetClaims(ctx *gin.Context) *Claims {
	return GetToken(ctx).Claims()
}
