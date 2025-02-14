package jwt

import (
	"errors"

	"github.com/gin-gonic/gin"
	jwt_ "github.com/golang-jwt/jwt/v5"
	"github.com/golang-jwt/jwt/v5/request"
)

var (
	ErrNoJWT         = errors.New("no jwt token")
	ErrUnequalSigner = errors.New("signer is not the same as token signer")
)

// Token represents a JSON Web Token with additional signed string representation.
type Token struct {
	*jwt_.Token
	signedToken string
}

// NewToken creates a new JWT token with the given configuration and claims.
func NewToken(config *Config, claims *Claims) (token *Token, err error) {
	token = &Token{
		Token: jwt_.NewWithClaims(config.Signer, claims),
	}
	token.signedToken, err = token.SignedString(config.getSecretKey())
	return
}

// NewTokenFromExtractor creates a new JWT token from an extracted token string and error.
// Intended for use with custom token extractors. For standard header and cookie extraction,
// use NewTokenFromHeader or NewTokenFromCookie instead.
func NewTokenFromExtractor(config *Config, extractorToken string, extractorErr error) (token *Token, err error) {
	if extractorErr != nil {
		return nil, extractorErr
	}
	if extractorToken == "" {
		return nil, ErrNoJWT
	}
	token = &Token{signedToken: extractorToken}
	token.Token, err = jwt_.ParseWithClaims(token.signedToken, &Claims{}, func(t *jwt_.Token) (interface{}, error) {
		if t.Method.Alg() != config.Signer.Alg() {
			return nil, ErrUnequalSigner
		}
		return config.getSecretKey(), nil
	})
	if err != nil {
		return nil, err
	}
	return
}

// NewTokenFromHeader creates a new JWT token by extracting it from the request header.
func NewTokenFromHeader(config *Config, ctx *gin.Context) (token *Token, err error) {
	extractorToken, extractorErr := request.BearerExtractor{}.ExtractToken(ctx.Request)
	return NewTokenFromExtractor(config, extractorToken, extractorErr)
}

// NewTokenFromCookie creates a new JWT token by extracting it from the request cookie.
func NewTokenFromCookie(config *Config, ctx *gin.Context) (token *Token, err error) {
	extractorToken, extractorErr := ctx.Cookie(config.CookieName)
	return NewTokenFromExtractor(config, extractorToken, extractorErr)
}

// NewTokenFromRequest creates a new JWT token by attempting to extract it from the request header first,
// and if not found, then from the request cookie.
func NewTokenFromRequest(config *Config, ctx *gin.Context) (token *Token, err error) {
	// Try to extract token from header first
	token, err = NewTokenFromHeader(config, ctx)

	// If no token in header, try to extract from cookie
	if errors.Is(err, request.ErrNoTokenInRequest) || errors.Is(err, ErrNoJWT) {
		token, err = NewTokenFromCookie(config, ctx)
	}
	return
}

// String returns the signed string representation of the JWT token.
func (r *Token) String() string {
	return r.signedToken
}

// Claims returns the custom claims associated with the JWT token.
// If the token is nil or the claims are not of type *Claims, it returns an empty Claims object.
func (r *Token) Claims() *Claims {
	if r.Token != nil {
		if claims, ok := r.Token.Claims.(*Claims); ok {
			return claims
		}
	}
	return &Claims{}
}
