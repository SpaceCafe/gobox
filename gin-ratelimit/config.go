package ratelimit

import (
	"errors"
	"time"
)

var (
	ErrInvalidMaxBurstRequests      = errors.New("max burst requests must be greater than 0")
	ErrInvalidMaxConcurrentRequests = errors.New("max concurrent requests must be greater than 0")
	ErrInvalidRequestQueueSize      = errors.New("request queue size must be greater than 0")
	ErrInvalidBurstDuration         = errors.New("burst duration must be greater than 0")
	ErrInvalidRequestTimeout        = errors.New("request timeout must be greater than 0")
)

// Config holds configuration related to rate limiting.
type Config struct {
	// MaxBurstRequests represents the maximum number of requests that can be processed in a burst.
	MaxBurstRequests int `json:"max_burst_requests" yaml:"max_burst_requests" mapstructure:"max_burst_requests"`

	// MaxConcurrentRequests represents the maximum number of concurrent processing slots available.
	MaxConcurrentRequests int `json:"max_concurrent_requests" yaml:"max_concurrent_requests" mapstructure:"max_concurrent_requests"`

	// RequestQueueSize represents the size of the waiting queue for incoming requests.
	RequestQueueSize int `json:"request_queue_size" yaml:"request_queue_size" mapstructure:"request_queue_size"`

	// BurstDuration represents the time span within which the burst limit is applied.
	BurstDuration time.Duration `json:"burst_duration" yaml:"burst_duration" mapstructure:"burst_duration"`

	// RequestTimeout represents the maximum duration a request can stay in the queue before being canceled.
	RequestTimeout time.Duration `json:"request_timeout" yaml:"request_timeout" mapstructure:"request_timeout"`
}

// NewConfig creates and returns a new Config with default values.
func NewConfig() *Config {
	return &Config{
		MaxBurstRequests:      20,
		MaxConcurrentRequests: 10,
		RequestQueueSize:      100,
		BurstDuration:         30 * time.Second,
		RequestTimeout:        1 * time.Minute,
	}
}

// Validate ensures that all necessary configurations are filled and within valid confines.
// Any misconfiguration results in well-defined standardized errors.
func (c *Config) Validate() error {
	if c.MaxBurstRequests <= 0 {
		return ErrInvalidMaxBurstRequests
	}
	if c.MaxConcurrentRequests <= 0 {
		return ErrInvalidMaxConcurrentRequests
	}
	if c.RequestQueueSize <= 0 {
		return ErrInvalidRequestQueueSize
	}
	if c.BurstDuration <= 0 {
		return ErrInvalidBurstDuration
	}
	if c.RequestTimeout <= 0 {
		return ErrInvalidRequestTimeout
	}
	return nil
}
