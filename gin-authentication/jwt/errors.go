package jwt

import (
	"errors"
)

var (
	ErrInvalidAudiences  = errors.New("jwt audiences must not be nil or empty")
	ErrInvalidCookieName = errors.New("jwt cookie names must not be empty or contain invalid characters")
	ErrInvalidIssuer     = errors.New("jwt issuer must not be empty")
	ErrInvalidSecret     = errors.New("jwt secrets must not be empty or shorter than 32 bytes")
	ErrInvalidSigner     = errors.New("jwt signer must not be nil")
	ErrInvalidTokenTTL   = errors.New("jwt tokens' TTL must be greater than zero")
	ErrSignerUnequal     = errors.New("jwt signer must be equal to the signer used to sign the token")
	ErrTokenIsNil        = errors.New("jwt token is nil")
	ErrClaimsIsNil       = errors.New("jwt claims is nil")
)
