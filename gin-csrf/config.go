package csrf

import (
	"crypto/sha256"
	"hash"
	"net/http"
	"regexp"

	"github.com/spacecafe/gobox/config"
)

var (
	_ config.Configure = (*Config)(nil)

	cookieNameValidator = regexp.MustCompile(`^[!#$%&'*\+\-.^_` + "`" + `|~0-9a-zA-Z]+$`)
	headerNameValidator = regexp.MustCompile(`^[A-Za-z0-9-]+$`)
)

// Config holds configuration related to CSRF protection.
type Config struct {
	// Secret is used to generate and validate CSRF tokens.
	Secret Secret `json:"secret" mapstructure:"secret" yaml:"secret"`

	// CookieName is the name of the cookie where the CSRF token will be stored.
	CookieName string `json:"cookieName" mapstructure:"cookie-name" yaml:"cookieName"`

	// CookieSameSite is the same site attribute of the cookie. Possible values are "Strict", "Lax" or "None".
	CookieSameSite CookieSameSite `json:"cookieSameSite" mapstructure:"cookie-same-site" yaml:"cookieSameSite"`

	// HeaderName is the name of the header where the CSRF token will be expected in requests.
	HeaderName string `json:"headerName" mapstructure:"header-name" yaml:"headerName"`

	// Signer is a function that returns a hash.Hash instance used creating an HMAC signature.
	Signer func() hash.Hash
}

// SetDefaults initializes the default values for the relevant fields in the struct.
func (r *Config) SetDefaults() {
	r.CookieName = "XSRF-TOKEN"
	r.HeaderName = "X-XSRF-TOKEN"
	r.Signer = sha256.New
}

// Validate ensures the all necessary configurations are filled and within valid confines.
func (r *Config) Validate() error {
	//nolint:mnd // Minimum length of 32 is required for HS256.
	if len(r.Secret) < 32 {
		return ErrInvalidSecret
	}

	if r.CookieName == "" || !cookieNameValidator.MatchString(r.CookieName) {
		return ErrInvalidCookieName
	}

	if r.CookieSameSite.SameSite < http.SameSiteLaxMode {
		return ErrInvalidCookieSameSite
	}

	if r.HeaderName == "" || !headerNameValidator.MatchString(r.HeaderName) {
		return ErrInvalidHeaderName
	}

	if r.Signer == nil {
		return ErrNoSigner
	}

	return nil
}
