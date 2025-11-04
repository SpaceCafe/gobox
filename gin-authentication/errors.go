package authentication

import (
	"errors"
)

var (
	ErrEmptyPassword         = errors.New("authentication password must not be empty")
	ErrEmptyToken            = errors.New("authentication token must not be empty")
	ErrInvalidAuthenticators = errors.New("authentication authenticators must not be nil")
	ErrInvalidMethod         = errors.New("authentication method is invalid")
	ErrInvalidPrincipals     = errors.New("authentication principals must not be nil")
	ErrInvalidRepository     = errors.New("authentication repository must not be nil")
	ErrInvalidTokens         = errors.New("authentication tokens must not be nil")
	ErrPrincipalNotFound     = errors.New("authentication principal not found")
	ErrSecretsNotEqual       = errors.New("authentication secrets are not equal")
)
