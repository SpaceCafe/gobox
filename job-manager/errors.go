package job_manager

import (
	"errors"
)

var (
	ErrInvalidBackend       = errors.New("job-manager backend is not valid")
	ErrInvalidPort          = errors.New("job-manager port must be a number between 1 and 65535")
	ErrInvalidRedisTTL      = errors.New("job-manager redis ttl must be greater than 0")
	ErrInvalidTimeout       = errors.New("job-manager timeout must be greater than 0")
	ErrJobManagerTerminated = errors.New("job-manager was terminated")
	ErrNoHost               = errors.New("job-manager host cannot be empty")
	ErrNoJobPointer         = errors.New("job must be a pointer")
	ErrNoRedisNamespace     = errors.New("job-manager redis namespace cannot be empty")
	ErrNoWorkerName         = errors.New("job-manager worker name cannot be empty")
	ErrTimeoutExceeded      = errors.New("job timeout exceeded")
)
