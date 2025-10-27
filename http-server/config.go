package http_server

import (
	"path"
	"strings"
	"time"
)

// Config defines the essential parameters for serving an http server.
type Config struct {

	// Host represents network host address.
	Host string `json:"host" yaml:"host" mapstructure:"host"`

	// BasePath represents the prefixed path in the URL.
	BasePath string `json:"base_path" yaml:"base_path" mapstructure:"base_path"`

	// CertFile represents the path to the certificate file.
	CertFile string `json:"cert_file" yaml:"cert_file" mapstructure:"cert_file"`

	// KeyFile represents the path to the key file.
	KeyFile string `json:"key_file" yaml:"key_file" mapstructure:"key_file"`

	// ReadTimeout represents the maximum duration before timing out read of the request.
	ReadTimeout time.Duration `json:"read_timeout" yaml:"read_timeout" mapstructure:"read_timeout"`

	// ReadHeaderTimeout represents the amount of time allowed to read request headers.
	ReadHeaderTimeout time.Duration `json:"read_header_timeout" yaml:"read_header_timeout" mapstructure:"read_header_timeout"`

	// Port specifies the port to be used for connections.
	Port int `json:"port" yaml:"port" mapstructure:"port"`
}

// SetDefaults initializes the default values for the relevant fields in the struct.
func (r *Config) SetDefaults() {
	r.Host = "127.0.0.1"
	r.ReadTimeout = time.Second * 30       //nolint:mnd // Default timeout value
	r.ReadHeaderTimeout = time.Second * 10 //nolint:mnd // Default header timeout value
	r.Port = 8080
}

// Validate ensures the all necessary configurations are filled and within valid confines.
func (r *Config) Validate() error {
	if r.Host == "" {
		return ErrNoHost
	}

	if r.BasePath != "" && (!path.IsAbs(r.BasePath) || strings.HasSuffix(r.BasePath, "/")) {
		return ErrInvalidBasePath
	}

	if r.CertFile != "" || r.KeyFile != "" {
		if r.CertFile == "" {
			return ErrNoCertFile
		}

		if r.KeyFile == "" {
			return ErrNoKeyFile
		}
	}

	if r.ReadTimeout <= 0 {
		return ErrInvalidReadTimeout
	}

	if r.ReadHeaderTimeout <= 0 {
		return ErrInvalidReadHeaderTimeout
	}

	if r.Port <= 0 || r.Port > 65535 {
		return ErrInvalidPort
	}

	return nil
}
