package terminator

import (
	"time"

	"github.com/spacecafe/gobox/config"
)

var _ config.Configure = (*Config)(nil)

// Config defines the essential parameters for serving the terminator.
type Config struct {
	// Timeout specifies the duration before the application is forcefully killed.
	Timeout time.Duration `json:"timeout" mapstructure:"timeout" yaml:"timeout"`

	// Force indicates whether to forcibly terminate the application without waiting for a graceful shutdown.
	Force bool `json:"force" mapstructure:"force" yaml:"force"`
}

// SetDefaults initializes the default values for the relevant fields in the struct.
func (r *Config) SetDefaults() {
	r.Timeout = time.Second * 3 //nolint:mnd // Default timeout value
	r.Force = true
}

// Validate ensures the all necessary configurations are filled and within valid confines.
func (r *Config) Validate() error {
	if r.Timeout <= 0 {
		return ErrInvalidTimeout
	}

	return nil
}
