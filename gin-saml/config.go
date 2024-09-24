package saml

import (
	"errors"
	"net/url"
	"slices"
	"strings"
	"time"
	"unicode"

	"github.com/crewjam/saml"
	"github.com/spacecafe/gobox/logger"
)

var (
	bindingURNs        = []string{saml.HTTPPostBinding, saml.HTTPRedirectBinding, saml.HTTPArtifactBinding, saml.SOAPBinding, saml.SOAPBindingV1}
	validSameSite      = []string{"strict", "lax", "none"}
	validNameIDFormats = []string{
		string(saml.UnspecifiedNameIDFormat),
		string(saml.EmailAddressNameIDFormat),
		string(saml.PersistentNameIDFormat),
		string(saml.TransientNameIDFormat),
	}

	ErrInvalidLogoutBindings = errors.New("logout bindings contains a not valid urn")
	ErrNoEntityID            = errors.New("entity id cannot be empty")
	ErrInvalidIDPMetadataURL = errors.New("idp metadata url is not valid")
	ErrNoCertFile            = errors.New("key file is set but cert_file is empty")
	ErrNoKeyFile             = errors.New("cert file is set but key_file is empty")
	ErrInvalidNameIDFormat   = errors.New("authn name id format is not valid")
	ErrInvalidURI            = errors.New("uri is not valid")
	ErrInvalidRedirectURI    = errors.New("default redirect uri is not valid")
	ErrInvalidErrorURI       = errors.New("default error uri is not valid")
	ErrInvalidCookieSameSite = errors.New("cookie same site is not valid")
	ErrInvalidCookieName     = errors.New("cookie name contains invalid characters or starts with '$'")
	ErrInvalidMaxIssueDelay  = errors.New("max issue delay must be greater than 0")
)

// Config holds configuration related to SAML as an authentication provider.
type Config struct {
	// LogoutBindings represents a list of bindings that can be used for logout requests.
	LogoutBindings []string `json:"logout_bindings" yaml:"logout_bindings" mapstructure:"logout_bindings"`

	// EntityID is the name of the service provider.
	EntityID string `json:"entity_id" yaml:"entity_id" mapstructure:"entity_id"`

	// IDPMetadataURL is the URL to the metadata configuration file of the identity provider.
	IDPMetadataURL string `json:"idp_metadata_url" yaml:"idp_metadata_url" mapstructure:"idp_metadata_url"`

	// CertFile represents the path to the certificate file.
	CertFile string `json:"cert_file" yaml:"cert_file" mapstructure:"cert_file"`

	// KeyFile represents the path to the key file.
	KeyFile string `json:"key_file" yaml:"key_file" mapstructure:"key_file"`

	// AuthnNameIDFormat is the format of the Name Identifier used in authentication requests.
	AuthnNameIDFormat string `json:"authn_name_id_format" yaml:"authn_name_id_format" mapstructure:"authn_name_id_format"`

	// URI represents the schema, domain and (optional) port of the service provider.
	URI string `json:"uri" yaml:"uri" mapstructure:"uri"`

	// DefaultRedirectURI is the default redirect URI used in authentication requests.
	DefaultRedirectURI string `json:"default_redirect_uri" yaml:"default_redirect_uri" mapstructure:"default_redirect_uri"`

	// DefaultErrorURI is the default error URI used in authentication requests.
	DefaultErrorURI string `json:"default_error_uri" yaml:"default_error_uri" mapstructure:"default_error_uri"`

	// CookieSameSite specifies the cookie SameSite attribute.
	CookieSameSite string `json:"cookie_same_site" yaml:"cookie_same_site" mapstructure:"cookie_same_site"`

	// CookieName is the name of the session cookie used for SAML authentication.
	CookieName string `json:"cookie_name" yaml:"cookie_name" mapstructure:"cookie_name"`

	// Logger specifies the used logger instance.
	Logger *logger.Logger

	// MaxIssueDelay is the maximum allowed delay for issuing SAML tokens.
	MaxIssueDelay time.Duration `json:"max_issue_delay" yaml:"max_issue_delay" mapstructure:"max_issue_delay"`

	// Mapping maps attributes from the identity provider to local attributes.
	Mapping map[string]string `json:"mapping" yaml:"mapping" mapstructure:"mapping"`

	// AllowIDPInitiated specifies whether IDP-initiated SAML authentication is allowed or not.
	AllowIDPInitiated bool `json:"allow_idp_initiated" yaml:"allow_idp_initiated" mapstructure:"allow_idp_initiated"`

	// SignRequest defines whether the requests are signed or not.
	SignRequest bool `json:"sign" yaml:"sign" mapstructure:"sign"`

	// UseArtifactResponse specifies whether to use artifact responses for authentication or not.
	UseArtifactResponse bool `json:"use_artifact_response" yaml:"use_artifact_response" mapstructure:"use_artifact_response"`

	// ForceAuthn forces the user to authenticate again even if they have a valid session.
	ForceAuthn bool `json:"force_authn" yaml:"force_authn" mapstructure:"force_authn"`
}

// NewConfig creates and returns a new Config having default values.
func NewConfig(log *logger.Logger) *Config {
	config := &Config{
		LogoutBindings:      []string{saml.HTTPPostBinding},
		AuthnNameIDFormat:   string(saml.TransientNameIDFormat),
		DefaultRedirectURI:  "/",
		DefaultErrorURI:     "/",
		CookieSameSite:      "strict",
		CookieName:          "token",
		MaxIssueDelay:       30 * time.Minute,
		AllowIDPInitiated:   false,
		SignRequest:         true,
		UseArtifactResponse: false,
		ForceAuthn:          false,
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
	for _, binding := range r.LogoutBindings {
		if !slices.Contains(bindingURNs, binding) {
			return ErrInvalidLogoutBindings
		}
	}

	if r.EntityID == "" {
		return ErrNoEntityID
	}

	if _, err := url.Parse(r.IDPMetadataURL); err != nil {
		return ErrInvalidIDPMetadataURL
	}

	if r.CertFile == "" {
		return ErrNoCertFile
	}

	if r.KeyFile == "" {
		return ErrNoKeyFile
	}

	if !slices.Contains(validNameIDFormats, r.AuthnNameIDFormat) {
		return ErrInvalidNameIDFormat
	}

	if _, err := url.ParseRequestURI(r.URI); err != nil {
		return ErrInvalidURI
	}

	if _, err := url.ParseRequestURI(r.DefaultRedirectURI); err != nil {
		return ErrInvalidRedirectURI
	}

	if _, err := url.ParseRequestURI(r.DefaultErrorURI); err != nil {
		return ErrInvalidErrorURI
	}

	if !slices.Contains(validSameSite, strings.ToLower(r.CookieSameSite)) {
		return ErrInvalidCookieSameSite
	}

	if !isValidCookieName(r.CookieName) {
		return ErrInvalidCookieName
	}

	if r.MaxIssueDelay <= 0 {
		return ErrInvalidMaxIssueDelay
	}

	return nil
}

// isValidCookieName checks if a given cookie name is valid.
// A valid cookie name must not be empty, must not start with a '$',
// and must only contain alphanumeric characters or underscores.
func isValidCookieName(name string) bool {
	if name == "" || name[0] == '$' {
		return false
	}
	for _, r := range name {
		if !(unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_') {
			return false
		}
	}
	return true
}
