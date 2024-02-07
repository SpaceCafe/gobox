package authentication

import "errors"

var (
	ErrNoPassword = errors.New("password of a user cannot be empty")
)

// Config holds configuration related to user and API key authentication.
type Config struct {

	// APIKeys list representing API keys that can be used to authenticate requests.
	APIKeys []string `json:"api_keys" yaml:"api_keys" mapstructure:"api_keys"`

	// HeaderName is the name of the header that contains the authentication token.
	// This could be "Authorization" or "API-Key", for example, which is commonly used in HTTP authentication.
	HeaderName string `json:"header_name" yaml:"header_name" mapstructure:"header_name"`

	// HeaderValuePrefix is a prefix that will be added to the API key in the header.
	// This can be used to provide additional context or information about how the API key was obtained,
	// such as "Bearer ", which is commonly used in HTTP authentication.
	HeaderValuePrefix string `json:"header_value_prefix" yaml:"header_value_prefix" mapstructure:"header_value_prefix"`

	// Users is a map where keys are usernames and values are passwords. It's used for basic HTTP authentication.
	Users map[string]string `json:"users" yaml:"users" mapstructure:"users"`
}

// NewConfig creates and returns a new Config having default values.
func NewConfig() *Config {
	return &Config{HeaderName: "API-Key"}
}

// Validate ensures the all necessary configurations are filled and within valid confines.
// Any misconfiguration results in well-defined standardized errors.
func (r *Config) Validate() error {
	for i := range r.Users {
		if len(r.Users[i]) == 0 {
			return ErrNoPassword
		}
	}

	return nil
}
