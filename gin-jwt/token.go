package jwt

import (
	"errors"

	"github.com/gin-gonic/gin"
	jwt_ "github.com/golang-jwt/jwt/v5"
	"github.com/golang-jwt/jwt/v5/request"
)

var (
	ErrUnequalSigner = errors.New("signer is not the same as token signer")
)

type Token struct {
	*jwt_.Token
	signedToken string
}

func NewToken(config *Config, claims *Claims) (*Token, error) {
	var err error
	token := &Token{
		Token: jwt_.NewWithClaims(config.Signer, claims),
	}
	token.signedToken, err = token.SignedString(config.getSecretKey())
	return token, err
}

func NewTokenFromHeader(config *Config, ctx *gin.Context) (*Token, error) {
	var err error
	token := &Token{}
	token.signedToken, err = request.BearerExtractor{}.ExtractToken(ctx.Request)
	if err != nil {
		return nil, err
	}

	token.Token, err = jwt_.ParseWithClaims(token.signedToken, &Claims{},
		func(t *jwt_.Token) (interface{}, error) {
			if t.Method.Alg() != config.Signer.Alg() {
				return nil, ErrUnequalSigner
			}

			return config.getSecretKey(), nil
		},
	)

	return token, err
}

func (r *Token) String() string {
	return r.signedToken
}
