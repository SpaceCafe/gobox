package httpserver

import (
	"path"
	"strings"
	"time"

	"github.com/spacecafe/gobox/config"
)

var _ config.Configure = (*Config)(nil)

// Config defines the essential parameters for serving an http server.
type Config struct {
	// Host represents network host address.
	Host string `json:"host" mapstructure:"host" yaml:"host"`

	// BasePath represents the prefixed path in the URL.
	BasePath string `json:"basePath" mapstructure:"base-path" yaml:"basePath"`

	// CertFile represents the path to the certificate file.
	CertFile string `json:"certFile" mapstructure:"cert-file" yaml:"certFile"`

	// KeyFile represents the path to the key file.
	KeyFile string `json:"keyFile" mapstructure:"key-file" yaml:"keyFile"`

	// ReadTimeout represents the maximum duration before timing out read of the request.
	ReadTimeout time.Duration `json:"readTimeout" mapstructure:"read-timeout" yaml:"readTimeout"`

	// ReadHeaderTimeout represents the amount of time allowed to read request headers.
	ReadHeaderTimeout time.Duration `json:"readHeaderTimeout" mapstructure:"read-header-timeout" yaml:"readHeaderTimeout"`

	// Port specifies the port to be used for connections.
	Port int `json:"port" mapstructure:"port" yaml:"port"`
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
