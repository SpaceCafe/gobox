package csrf

import (
	"crypto/sha256"
	"errors"
	"hash"
	"net/http"
	"path"
	"regexp"
	"strings"

	"github.com/spacecafe/gobox/logger"
)

const (
	DefaultCookieName     = "csrf_token"
	DefaultHeaderName     = "X-CSRF-Token"
	DefaultSameSite       = "strict"
	DefaultPath           = "/"
	DefaultDomain         = ""
	DefaultTokenLength    = 32
	DefaultHTTPOnlyCookie = true
	DefaultSecureCookie   = true
)

var (
	ErrInvalidExcludedRoutes = errors.New("excluded routes must not be nil")
	ErrInvalidExcludedRoute  = errors.New("excluded route must be absolute and not end with a slash")
	ErrNoSecretKey           = errors.New("secret key must not be empty")
	ErrInvalidCookieName     = errors.New("cookie name must not be empty or contain invalid characters")
	ErrInvalidHeaderName     = errors.New("header name must not be empty or contain invalid characters")
	ErrInvalidSameSite       = errors.New("same site is not valid")
	ErrInvalidPath           = errors.New("path must be absolute and not end with a slash")
	ErrInvalidTokenLength    = errors.New("token length must be greater than 0")
	ErrNoSigner              = errors.New("signer must not be nil")
	ErrNoLogger              = errors.New("logger cannot be empty")

	//nolint:gochecknoglobals // Maintain a set of predefined http.SameSite that are used throughout the application.
	SameSites = map[string]http.SameSite{
		"lax":    http.SameSiteLaxMode,
		"strict": http.SameSiteStrictMode,
		"none":   http.SameSiteNoneMode,
	}

	validCookieName = regexp.MustCompile(`^[!#$%&'*\+\-.^_` + "`" + `|~0-9a-zA-Z]+$`)
	validHeaderName = regexp.MustCompile(`^[A-Za-z0-9-]+$`)
)

// Config holds configuration related to CSRF protection.
type Config struct {

	// ExcludedRoutes is a list of routes that are excluded from CSRF protection.
	ExcludedRoutes []string `json:"excluded_routes" yaml:"excluded_routes" mapstructure:"excluded_routes"`

	// SecretKey is used to generate and validate CSRF tokens.
	SecretKey []byte `json:"secret_key" yaml:"secret_key" mapstructure:"secret_key"`

	// CookieName is the name of the cookie where the CSRF token will be stored.
	CookieName string `json:"cookie_name" yaml:"cookie_name" mapstructure:"cookie_name"`

	// HeaderName is the name of the header where the CSRF token will be expected in requests.
	HeaderName string `json:"header_name" yaml:"header_name" mapstructure:"header_name"`

	// SameSite is the SameSite policy for the CSRF cookie.
	SameSite string `json:"same_site" yaml:"same_site" mapstructure:"same_site"`

	// Path is the path for which the CSRF cookie is valid.
	Path string `json:"path" yaml:"path" mapstructure:"path"`

	// Domain is the domain for which the CSRF cookie is valid.
	Domain string `json:"domain" yaml:"domain" mapstructure:"domain"`

	// TokenLength is the length of the CSRF token.
	TokenLength int `json:"token_length" yaml:"token_length" mapstructure:"token_length"`

	// Signer is a function that returns a new hash.Hash to be used for signing CSRF tokens.
	Signer func() hash.Hash

	// sameSite is the internal representation of SameSite.
	sameSite http.SameSite

	// Logger specifies the used logger instance.
	Logger *logger.Logger

	// HTTPOnlyCookie indicates whether the CSRF cookie should be marked as HTTP only.
	HTTPOnlyCookie bool `json:"http_only_cookie" yaml:"http_only_cookie" mapstructure:"http_only_cookie"`

	// SecureCookie indicates whether the CSRF cookie should be marked as secure (HTTPS only).
	SecureCookie bool `json:"secure_cookie" yaml:"secure_cookie" mapstructure:"secure_cookie"`
}

// NewConfig creates and returns a new Config having default values.
func NewConfig(log *logger.Logger) *Config {
	config := &Config{
		ExcludedRoutes: make([]string, 0),
		CookieName:     DefaultCookieName,
		HeaderName:     DefaultHeaderName,
		SameSite:       DefaultSameSite,
		Path:           DefaultPath,
		Domain:         DefaultDomain,
		TokenLength:    DefaultTokenLength,
		Signer:         sha256.New,
		HTTPOnlyCookie: DefaultHTTPOnlyCookie,
		SecureCookie:   DefaultSecureCookie,
	}

	if log != nil {
		config.Logger = log
	} else {
		config.Logger = logger.Default()
	}

	return config
}

// Validate ensures the all necessary configurations are filled and within valid confines.
// Any misconfiguration results in well-defined standardized errors.
func (r *Config) Validate() error {
	if r.ExcludedRoutes == nil {
		return ErrInvalidExcludedRoutes
	}
	for i := range r.ExcludedRoutes {
		if !path.IsAbs(r.ExcludedRoutes[i]) {
			return ErrInvalidExcludedRoute
		}
	}
	if len(r.SecretKey) == 0 {
		return ErrNoSecretKey
	}
	if validCookieName.MatchString(r.CookieName) {
		return ErrInvalidCookieName
	}
	if validHeaderName.MatchString(r.HeaderName) {
		return ErrInvalidHeaderName
	}
	if err := r.SetSameSite(""); err != nil {
		return err
	}
	if !path.IsAbs(r.Path) {
		return ErrInvalidPath
	}
	if r.TokenLength <= 0 {
		return ErrInvalidTokenLength
	}
	if r.Signer == nil {
		return ErrNoSigner
	}
	if r.Logger == nil {
		return ErrNoLogger
	}
	return nil
}

// SetSameSite sets the SameSite attribute for the Config struct.
// It takes a string parameter SameSite and returns an error if the value is invalid.
// If SameSite is an empty string, it uses the existing value in the Config struct.
func (r *Config) SetSameSite(sameSite string) error {
	var key = sameSite
	if sameSite == "" {
		key = r.SameSite
	}
	if value, ok := SameSites[strings.ToLower(key)]; ok {
		r.sameSite = value
		return nil
	}
	return ErrInvalidSameSite
}
