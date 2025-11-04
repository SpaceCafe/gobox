package authentication

import (
	"github.com/spacecafe/gobox/config"
	"github.com/spacecafe/gobox/gin-authentication/jwt"
)

var (
	_ config.Configure = (*Config)(nil)
)

// Config holds configuration related to user and API key authentication.
type Config struct {
	// Tokens list representing API keys that can be used to authenticate requests.
	Tokens []string `json:"tokens" yaml:"tokens" mapstructure:"tokens"`

	// Authenticators is a list of authenticators that can be used to authenticate requests.
	Authenticators []Authenticator

	// Repository is the repository used to store and retrieve user information.
	Repository Repository

	// Principals is a map of principal ids to passwords.
	Principals map[string]string `json:"principals" yaml:"principals" mapstructure:"principals"`

	JWT *jwt.Config `json:"jwt" yaml:"jwt" mapstructure:"jwt"`
}

// SetDefaults initializes the default values for the relevant fields in the struct.
func (r *Config) SetDefaults() {
	r.Tokens = []string{}
	r.Authenticators = []Authenticator{
		NewTokenAuthenticator(r),
	}
	r.Repository = NewConfigRepository(r)
	r.Principals = map[string]string{}
	r.JWT = &jwt.Config{}
	r.JWT.SetDefaults()
}

// Validate ensures the all necessary configurations are filled and within valid confines.
func (r *Config) Validate() error {
	if r.Tokens == nil {
		return ErrInvalidTokens
	}
	for i := range r.Tokens {
		if r.Tokens[i] == "" {
			return ErrEmptyToken
		}
	}

	if r.Authenticators == nil {
		return ErrInvalidAuthenticators
	}
	for i := range r.Authenticators {
		if r.Authenticators[i] == nil {
			return ErrInvalidAuthenticators
		}
		switch r.Authenticators[i].(type) {
		case *BearerAuthenticator, *JWTAuthenticator:
			if err := r.JWT.Validate(); err != nil {
				return err
			}
		}
	}

	if r.Repository == nil {
		return ErrInvalidRepository
	}

	if r.Principals == nil {
		return ErrInvalidPrincipals
	}
	for i := range r.Principals {
		if r.Principals[i] == "" {
			return ErrEmptyPassword
		}
	}

	return nil
}
