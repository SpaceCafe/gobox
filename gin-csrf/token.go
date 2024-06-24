package csrf

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"hash"

	"github.com/gin-gonic/gin"
)

var (
	ErrTokenGeneration             = errors.New("failed to generate CSRF token")
	ErrCookieRetrieval             = errors.New("failed to retrieve CSRF cookie")
	ErrTokenDecoding               = errors.New("failed to decode CSRF token")
	ErrInvalidSubmittedTokenLength = errors.New("submitted CSRF token length is invalid")
	ErrInvalidTokenSignature       = errors.New("token signature is invalid")
)

// Token represents a signed cookie-based token with a randomly generated message,
// its signature, and a URL-friendly string.
type Token struct {
	Message      []byte
	Signature    []byte
	encodedToken string
	signer       hash.Hash
}

// NewToken generates a new CSRF token using the provided configuration.
func NewToken(config *Config) (*Token, error) {
	token := &Token{
		signer:  hmac.New(config.Signer, config.SecretKey),
		Message: make([]byte, config.TokenLength),
	}

	// Generate new message.
	_, err := rand.Read(token.Message)
	if err != nil {
		return nil, ErrTokenGeneration
	}

	// Sign new message
	token.signer.Reset()
	token.signer.Write(token.Message)
	token.Signature = token.signer.Sum(nil)

	// Encode token with base64
	token.encodedToken = base64.RawURLEncoding.EncodeToString(append(token.Message, token.Signature...))

	return token, nil
}

// NewTokenFromCookie retrieves and validates a CSRF token from a cookie.
func NewTokenFromCookie(config *Config, ctx *gin.Context) (*Token, error) {
	token := &Token{
		signer: hmac.New(config.Signer, config.SecretKey),
	}

	// Retrieve cookie
	var err error
	token.encodedToken, err = ctx.Cookie(config.CookieName)
	if err != nil {
		return nil, ErrCookieRetrieval
	}

	// Decode token with base64
	decodedToken, err := base64.RawURLEncoding.DecodeString(token.encodedToken)
	if err != nil {
		return nil, ErrTokenDecoding
	}

	// Check length of submitted token
	if len(decodedToken) <= token.signer.Size() {
		return nil, ErrInvalidSubmittedTokenLength
	}

	token.Message = decodedToken[:len(decodedToken)-token.signer.Size()]
	token.Signature = decodedToken[len(decodedToken)-token.signer.Size():]

	// Sign message
	token.signer.Reset()
	token.signer.Write(token.Message)

	if !hmac.Equal(token.Signature, token.signer.Sum(nil)) {
		return nil, ErrInvalidTokenSignature
	}
	return token, nil
}

// String returns the encoded token as a string.
func (r *Token) String() string {
	return r.encodedToken
}

// Compare securely compares the provided token with the stored token.
func (r *Token) Compare(token string) bool {
	if r.encodedToken == "" || len(token) != len(r.encodedToken) {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(token), []byte(r.encodedToken)) == 1
}
