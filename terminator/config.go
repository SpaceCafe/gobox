package terminator

import (
	"errors"
	"time"
)

const (
	DefaultTimeout = 3 * time.Second
)

var (
	ErrInvalidTimeout = errors.New("timeout must be greater than 0")
)

// Config holds configuration related to Terminator.
type Config struct {

	// Timeout specifies the duration before the application is forced killed.
	Timeout time.Duration
}

// NewConfig creates and returns a new Config having default values.
func NewConfig() *Config {
	config := &Config{
		Timeout: DefaultTimeout,
	}

	return config
}

// Validate ensures the all necessary configurations are filled and within valid confines.
// Any misconfiguration results in well-defined standardized errors.
func (r *Config) Validate() error {
	if r.Timeout <= 0 {
		return ErrInvalidTimeout
	}
	return nil
}
