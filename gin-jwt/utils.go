package jwt

import (
	"github.com/gin-gonic/gin"
)

func GetClaims(ctx *gin.Context) *Claims {
	if token, ok := ctx.Get("jwt/token"); ok {
		return token.(*Token).Claims.(*Claims)
	}
	return &Claims{}
}
