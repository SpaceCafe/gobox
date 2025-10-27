package job_manager

import (
	"os"
	"time"
)

// Config holds configuration related to the job manager.
type Config struct {
	// WorkerName is the name of the worker, defaulting to the hostname.
	WorkerName string `json:"worker_name" yaml:"worker_name" mapstructure:"worker_name"`

	// Backend specifies the type of backend to be used.
	Backend string `json:"backend" yaml:"backend" mapstructure:"backend"`

	// RedisHost is the hostname or IP address of the Redis server.
	RedisHost string `json:"redis_host" yaml:"redis_host" mapstructure:"redis_host"`

	// RedisUsername is the username for authenticating with the Redis server (optional).
	RedisUsername string `json:"redis_username" yaml:"redis_username" mapstructure:"redis_username"`

	// RedisPassword is the password for authenticating with the Redis server (optional).
	RedisPassword string `json:"redis_password" yaml:"redis_password" mapstructure:"redis_password"`

	// RedisNamespace is the namespace used to prefix keys in Redis.
	RedisNamespace string `json:"redis_namespace" yaml:"redis_namespace" mapstructure:"redis_namespace"`

	// RedisPort is the port number on which the Redis server is listening to.
	RedisPort int `json:"redis_port" yaml:"redis_port" mapstructure:"redis_port"`

	// RedisTTL is the time-to-live in seconds for a Redis key. This value determines how long the data should be cached before it expires.
	RedisTTL time.Duration `json:"redis_ttl" yaml:"redis_ttl" mapstructure:"redis_ttl"`

	// Timeout represents the amount of time allowed to wait for a job.
	Timeout time.Duration `json:"timeout" yaml:"timeout" mapstructure:"timeout"`
}

// SetDefaults initializes the default values for the relevant fields in the struct.
func (r *Config) SetDefaults() {
	r.WorkerName, _ = os.Hostname()
	r.Backend = BackendRedis
	r.RedisNamespace = "jobs"
	r.RedisPort = 6379
	r.RedisTTL = time.Hour
	r.Timeout = time.Second * 30 //nolint:mnd // Default timeout value
}

// Validate ensures the all necessary configurations are filled and within valid confines.
func (r *Config) Validate() error {
	if r.WorkerName == "" {
		return ErrNoWorkerName
	}

	switch r.Backend {
	case BackendRedis:
		if r.RedisHost == "" {
			return ErrNoHost
		}

		if r.RedisNamespace == "" {
			return ErrNoRedisNamespace
		}

		if r.RedisPort <= 0 || r.RedisPort > 65535 {
			return ErrInvalidPort
		}

		if r.RedisTTL <= 0 {
			return ErrInvalidRedisTTL
		}
	default:
		return ErrInvalidBackend
	}

	if r.Timeout <= 0 {
		return ErrInvalidTimeout
	}

	return nil
}
