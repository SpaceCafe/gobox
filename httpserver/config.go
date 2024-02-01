package httpserver

import (
	"errors"
	"path"
	"strings"
	"time"
)

const (
	DefaultHost              = "127.0.0.1"
	DefaultBasePath          = "/"
	DefaultReadTimeout       = time.Second * 30
	DefaultReadHeaderTimeout = time.Second * 10
	DefaultPort              = 8080
)

var (
	ErrNoHost                   = errors.New("host cannot be empty")
	ErrNoBasePath               = errors.New("base path must be an absolute path")
	ErrInvalidBasePath          = errors.New("base path must end with a trailing slash")
	ErrNoCertFile               = errors.New("key_file is set but cert_file is empty")
	ErrNoKeyFile                = errors.New("cert_file is set but key_file is empty")
	ErrInvalidReadTimeout       = errors.New("read_timeout must be greater than 0")
	ErrInvalidReadHeaderTimeout = errors.New("read_header_timeout must be greater than 0")
	ErrInvalidPort              = errors.New("port must be a number between 1 and 65535")
)

// Config defines the essential parameters for serving a Lambda broker service.
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

// NewConfig creates and returns a new Config having default values.
func NewConfig() *Config {
	return &Config{
		Host:              DefaultHost,
		BasePath:          DefaultBasePath,
		ReadTimeout:       DefaultReadTimeout,
		ReadHeaderTimeout: DefaultReadHeaderTimeout,
		Port:              DefaultPort,
	}
}

// Validate ensures the all necessary configurations are filled and within valid confines.
// This includes checks for host, certificates, port, and timeouts.
// Any misconfiguration results in well-defined standardized errors.
func (m *Config) Validate() error {
	if m.Host == "" {
		return ErrNoHost
	}

	if !path.IsAbs(m.BasePath) {
		return ErrNoBasePath
	}

	if !strings.HasSuffix(m.BasePath, "/") {
		return ErrInvalidBasePath
	}

	if m.CertFile != "" || m.KeyFile != "" {
		if m.CertFile == "" {
			return ErrNoCertFile
		}

		if m.KeyFile == "" {
			return ErrNoKeyFile
		}
	}

	if m.ReadTimeout <= 0 {
		return ErrInvalidReadTimeout
	}

	if m.ReadHeaderTimeout <= 0 {
		return ErrInvalidReadHeaderTimeout
	}

	if m.Port <= 0 || m.Port > 65535 {
		return ErrInvalidPort
	}

	return nil
}
