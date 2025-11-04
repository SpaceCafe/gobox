package jwt

import (
	"regexp"
	"time"

	jwt2 "github.com/golang-jwt/jwt/v5"
	"github.com/spacecafe/gobox/config"
)

var (
	_ config.Configure = (*Config)(nil)

	cookieNameValidator = regexp.MustCompile(`^[!#$%&'*\+\-.^_` + "`" + `|~0-9a-zA-Z]+$`)
)

// Config holds configuration related to JWT.
type Config struct {
	// Secret as a base64 encoded string (RFC 4648) is used to generate and validate access tokens.
	Secret Secret `json:"secret" yaml:"secret" mapstructure:"secret"`

	// RefreshSecret as a base64 encoded string (RFC 4648) is used to generate and validate refresh tokens.
	RefreshSecret Secret `json:"refresh_secret" yaml:"refresh_secret" mapstructure:"refresh_secret"`

	// Audience is the intended recipient of the token.
	// It is usually a list of URLs of the services that can consume the token.
	// Example: https://api.example.com
	Audience []string `json:"audience" yaml:"audience" mapstructure:"audience"`

	// Issuer is the entity that issues the token.
	// It is usually the URL of the authorization server.
	// Example: https://auth.example.com
	Issuer string `json:"issuer" yaml:"issuer" mapstructure:"issuer"`

	// CookieName is the name of the cookie that stores the access token.
	CookieName string `json:"cookie_name" yaml:"cookie_name" mapstructure:"cookie_name"`

	// RefreshCookieName is the name of the cookie that stores the refresh token.
	RefreshCookieName string `json:"refresh_cookie_name" yaml:"refresh_cookie_name" mapstructure:"refresh_cookie_name"`

	// Signer is a function that returns a new SigningMethod to be used for signing JWT.
	Signer jwt2.SigningMethod

	// AccessTokenTTL is the duration for which the access token is valid.
	AccessTokenTTL time.Duration `json:"access_token_ttl" yaml:"access_token_ttl" mapstructure:"access_token_ttl"`

	// RefreshTokenTTL is the duration for which the refresh token is valid.
	RefreshTokenTTL time.Duration `json:"refresh_token_ttl" yaml:"refresh_token_ttl" mapstructure:"refresh_token_ttl"`
}

// SetDefaults initializes the default values for the relevant fields in the struct.
func (r *Config) SetDefaults() {
	r.CookieName = "__Host-access_token"
	r.RefreshCookieName = "__Host-refresh_token"
	r.Signer = jwt2.SigningMethodHS256
	r.AccessTokenTTL = time.Hour * 15  //nolint:mnd // Default access token expiration time
	r.RefreshTokenTTL = time.Hour * 24 //nolint:mnd // Default refresh token expiration time
}

// Validate ensures the all necessary configurations are filled and within valid confines.
func (r *Config) Validate() error {
	//nolint:mnd // Minimum length of 32 is required for HS256.
	if len(r.Secret) < 32 {
		return ErrInvalidSecret
	}

	//nolint:mnd // Minimum length of 32 is required for HS256.
	if len(r.RefreshSecret) < 32 {
		return ErrInvalidSecret
	}

	if len(r.Audience) == 0 {
		return ErrInvalidAudiences
	}
	for i := range r.Audience {
		if r.Audience[i] == "" {
			return ErrInvalidAudiences
		}
	}

	if r.Issuer == "" {
		return ErrInvalidIssuer
	}

	if r.CookieName == "" || !cookieNameValidator.MatchString(r.CookieName) {
		return ErrInvalidCookieName
	}

	if r.RefreshCookieName == "" || !cookieNameValidator.MatchString(r.RefreshCookieName) {
		return ErrInvalidCookieName
	}

	if r.Signer == nil {
		return ErrInvalidSigner
	}

	if r.AccessTokenTTL == 0 {
		return ErrInvalidTokenTTL
	}

	if r.RefreshTokenTTL == 0 {
		return ErrInvalidTokenTTL
	}

	return nil
}
