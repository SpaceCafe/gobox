package jwt

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	jwt2 "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

// Token represents a JSON Web Token with additional signed string representation.
type Token struct {
	*jwt2.Token

	cfg       *Config
	tokenType TokenType
}

// New creates a new Token with the given claims.
func New(cfg *Config, claims *Claims, tokenType TokenType) *Token {
	token := &Token{
		cfg:       cfg,
		tokenType: tokenType,
	}

	claims.RegisteredClaims.ID = uuid.New().String()
	claims.Issuer = cfg.Issuer
	claims.Audience = cfg.Audience
	claims.RegisteredClaims.ExpiresAt = jwt2.NewNumericDate(time.Now().Add(token.TTL()))
	claims.NotBefore = jwt2.NewNumericDate(time.Now())
	claims.IssuedAt = jwt2.NewNumericDate(time.Now())
	token.Token = jwt2.NewWithClaims(token.cfg.Signer, claims)

	return token
}

// NewFromString creates a new Token from a signed string representation.
func NewFromString(cfg *Config, signedToken string, tokenType TokenType) (*Token, error) {
	var err error

	token := &Token{
		cfg:       cfg,
		tokenType: tokenType,
	}
	secret := token.secret()

	token.Token, err = jwt2.ParseWithClaims(
		signedToken,
		&Claims{},
		func(t *jwt2.Token) (any, error) {
			if t.Method.Alg() != cfg.Signer.Alg() {
				return nil, ErrSignerUnequal
			}

			return secret, nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	return token, nil
}

// Claims returns the claims of the token.
func (r *Token) Claims() *Claims {
	if r.Token != nil {
		if claims, ok := r.Token.Claims.(*Claims); ok {
			return claims
		}
	}

	return nil
}

// Cookie generates an HTTP cookie from the token with specific security constraints.
func (r *Token) Cookie() (*http.Cookie, error) {
	claims := r.Claims()
	if claims == nil {
		return nil, ErrNoClaims
	}

	value, err := r.SignedString()
	if err != nil {
		return nil, err
	}

	// Create the cookie with the following constraints:
	// - `Domain` is not set, so the cookie is a host-only cookie.
	// - `Expires` is not set, so the cookie is a session cookie.
	// - `Secure` is set, so the cookie is only sent over HTTPS.
	// - `HttpOnly` is set, so the cookie is inaccessible to client-side JavaScript.
	// - `Partitioned` is set, so the cookie is sent in a separate cookie jar.
	return &http.Cookie{
		Name:        r.cookieName(),
		Value:       url.QueryEscape(value),
		Path:        "/",
		MaxAge:      int(claims.ExpiresAt().Unix() - time.Now().Unix()),
		Secure:      true,
		HttpOnly:    true,
		SameSite:    r.cfg.CookieSameSite.SameSite,
		Partitioned: true,
	}, nil
}

// Renew updates the token's expiration time based on its type and re-signs it.
func (r *Token) Renew() error {
	claims := r.Claims()
	if claims == nil {
		return ErrNoClaims
	}

	claims.RegisteredClaims.ID = uuid.New().String()
	claims.RegisteredClaims.ExpiresAt = jwt2.NewNumericDate(time.Now().Add(r.TTL()))
	r.Token = jwt2.NewWithClaims(r.cfg.Signer, claims)

	return nil
}

// SignedString returns the signed string representation of the token.
func (r *Token) SignedString() (string, error) {
	if r.Token == nil {
		return "", ErrNoToken
	}

	//nolint:wrapcheck // wrap check is not relevant here.
	return r.Token.SignedString(r.secret())
}

// TTL returns the time-to-live duration for the token based on its type.
func (r *Token) TTL() time.Duration {
	if r.tokenType == RefreshToken {
		return r.cfg.RefreshTokenTTL
	}

	return r.cfg.AccessTokenTTL
}

// cookieName determines the name of the cookie based on the token type.
func (r *Token) cookieName() string {
	if r.tokenType == RefreshToken {
		return r.cfg.RefreshCookieName
	}

	return r.cfg.CookieName
}

// secret returns the secret bytes used for signing the token based on its type.
func (r *Token) secret() []byte {
	if r.tokenType == RefreshToken {
		return r.cfg.RefreshSecret.Bytes()
	}

	return r.cfg.Secret.Bytes()
}
