package job_manager

import (
	"errors"
	"os"
	"slices"
	"time"

	"github.com/spacecafe/gobox/logger"
)

const (
	DefaultBackend        = "redis"
	DefaultRedisNamespace = "jobs"
	DefaultRedisPort      = 6379
	DefaultRedisTTL       = time.Hour
	DefaultTimeout        = time.Second * 30
)

var (
	validBackends = []string{"redis"}

	ErrNoWorkerName     = errors.New("worker name cannot be empty")
	ErrInvalidBackend   = errors.New("backend is not valid")
	ErrNoHost           = errors.New("host cannot be empty")
	ErrInvalidPort      = errors.New("port must be a number between 1 and 65535")
	ErrNoRedisNamespace = errors.New("redis namespace cannot be empty")
	ErrInvalidRedisTTL  = errors.New("redis ttl must be greater than 0")
	ErrInvalidTimeout   = errors.New("timeout must be greater than 0")
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

	// RedisPort is the port number on which the Redis server is listening.
	RedisPort int `json:"redis_port" yaml:"redis_port" mapstructure:"redis_port"`

	RedisTTL time.Duration `json:"redis_ttl" yaml:"redis_ttl" mapstructure:"redis_ttl"`

	// Timeout represents the amount of time allowed to wait for a job.
	Timeout time.Duration `json:"timeout" yaml:"timeout" mapstructure:"timeout"`

	// Logger specifies the used logger instance.
	Logger *logger.Logger
}

// NewConfig creates and returns a new Config having default values.
func NewConfig(log *logger.Logger) *Config {
	config := &Config{
		Backend:        DefaultBackend,
		RedisNamespace: DefaultRedisNamespace,
		RedisPort:      DefaultRedisPort,
		RedisTTL:       DefaultRedisTTL,
		Timeout:        DefaultTimeout,
	}
	config.WorkerName, _ = os.Hostname()

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
	if r.WorkerName == "" {
		return ErrNoWorkerName
	}

	if !slices.Contains(validBackends, r.Backend) {
		return ErrInvalidBackend
	}

	switch r.Backend {
	case "redis":
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
	}

	if r.Timeout <= 0 {
		return ErrInvalidTimeout
	}

	return nil
}
