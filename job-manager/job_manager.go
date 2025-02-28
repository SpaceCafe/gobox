package job_manager

import (
	"reflect"
)

// NewJobManager initializes a new job manager based on the provided configuration.
// It returns an instance of IJobManager or an error if the backend is invalid.
func NewJobManager[T any](config *Config) (jobManager IJobManager, err error) {
	if reflect.TypeFor[T]().Kind() != reflect.Struct {
		panic("model must be a struct")
	}
	if _, ok := any((*T)(nil)).(IJob); !ok {
		panic("model must implement IJob interface")
	}

	switch config.Backend {
	case "redis":
		return NewRedisJobManager[T](config)
	default:
		return nil, ErrInvalidBackend
	}
}
