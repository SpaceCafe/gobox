package jwt

import (
	"encoding/base64"
	"errors"
	"path"
	"time"

	jwt_ "github.com/golang-jwt/jwt/v5"
	"github.com/spacecafe/gobox/logger"
)

const (
	DefaultTokenExpiration = time.Hour * 24 // 1 day
	DefaultIssuer          = "API server"
)

var (
	DefaultExcludedRoutes = make([]string, 0)
	DefaultAudience       = []string{"api"}
	DefaultSigner         = jwt_.SigningMethodHS256

	ErrInvalidExcludedRoutes  = errors.New("excluded routes must not be nil")
	ErrInvalidExcludedRoute   = errors.New("excluded route must be absolute and not end with a slash")
	ErrNoSecretKey            = errors.New("secret key must not be empty")
	ErrInvalidAudiences       = errors.New("audiences must not be nil")
	ErrNoAudience             = errors.New("audience cannot be empty")
	ErrNoIssuer               = errors.New("issuer cannot be empty")
	ErrNoSigner               = errors.New("signer must not be nil")
	ErrInvalidTokenExpiration = errors.New("token expiration must be greater than zero")
	ErrNoLogger               = errors.New("logger cannot be empty")
)

// Config holds configuration related to JWT.
type Config struct {

	// ExcludedRoutes is a list of routes that are excluded from JWT.
	ExcludedRoutes []string `json:"excluded_routes" yaml:"excluded_routes" mapstructure:"excluded_routes"`

	// Audiences is the intended recipient of the token.
	Audiences []string `json:"audiences" yaml:"audiences" mapstructure:"audiences"`

	secretKey []byte

	// SecretKey as a base64 encoded string (RFC 4648) is used to generate and validate JWT.
	SecretKey string `json:"secret_key" yaml:"secret_key" mapstructure:"secret_key"`

	// Issuer is the entity that issues the token.
	Issuer string `json:"issuer" yaml:"issuer" mapstructure:"issuer"`

	// Signer is a function that returns a new SigningMethod to be used for signing JWT.
	Signer jwt_.SigningMethod

	// TokenExpiration is the duration for which the token is valid.
	TokenExpiration time.Duration `json:"token_expiration" yaml:"token_expiration" mapstructure:"token_expiration"`

	// Logger specifies the used logger instance.
	Logger *logger.Logger
}

// NewConfig creates and returns a new Config having default values.
func NewConfig(log *logger.Logger) *Config {
	config := &Config{
		ExcludedRoutes:  DefaultExcludedRoutes,
		Audiences:       DefaultAudience,
		Issuer:          DefaultIssuer,
		Signer:          DefaultSigner,
		TokenExpiration: DefaultTokenExpiration,
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
	if r.SecretKey == "" {
		return ErrNoSecretKey
	}
	if err := r.setSecretKey(r.SecretKey); err != nil {
		return err
	}
	if r.Audiences == nil {
		return ErrInvalidAudiences
	}
	for i := range r.Audiences {
		if len(r.Audiences[i]) == 0 {
			return ErrNoAudience
		}
	}
	if r.Issuer == "" {
		return ErrNoIssuer
	}
	if r.Signer == nil {
		return ErrNoSigner
	}
	if r.TokenExpiration == 0 {
		return ErrInvalidTokenExpiration
	}
	if r.Logger == nil {
		return ErrNoLogger
	}
	return nil
}

func (r *Config) setSecretKey(key string) (err error) {
	r.secretKey, err = base64.StdEncoding.DecodeString(key)
	return
}

func (r *Config) getSecretKey() []byte {
	if r.secretKey == nil || len(r.secretKey) == 0 {
		panic(ErrNoSecretKey)
	}
	return r.secretKey
}
