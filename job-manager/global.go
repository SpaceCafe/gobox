package job_manager

import (
	"reflect"

	"github.com/spacecafe/gobox/logger"
)

// New initializes a new job-manager based on the provided configuration.
// It returns an instance of Manager or an error if the backend is invalid.
func New[T any](cfg *Config, log logger.Logger) (Manager, error) {
	if reflect.TypeFor[T]().Kind() != reflect.Struct {
		panic("job-manager model must be a struct")
	}
	if _, ok := any((*T)(nil)).(Job); !ok {
		panic("job-manager model must implement Job interface")
	}

	switch cfg.Backend {
	case BackendRedis:
		return NewRedisManager[T](cfg, log)
	default:
		return nil, ErrInvalidBackend
	}
}
