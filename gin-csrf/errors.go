package csrf

import (
	"errors"
)

var (
	ErrInvalidCookieName = errors.New(
		"csrf cookie name must not be empty or contain invalid characters",
	)
	ErrInvalidCookieSameSite = errors.New("csrf cookie same site mode is invalid")
	ErrInvalidHeaderName     = errors.New(
		"csrf header name must not be empty or contain invalid characters",
	)
	ErrInvalidSecret    = errors.New("csrf secret must not be empty or shorter than 32 bytes")
	ErrInvalidSessionID = errors.New(
		"csrf session id must not be empty or shorter than 128 bits",
	)
	ErrInvalidToken = errors.New("csrf token is invalid")
	ErrNoSession    = errors.New("csrf session is missing")
	ErrNoSigner     = errors.New("csrf signer must not be nil")
)
