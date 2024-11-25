package job_manager

// NewJobManager initializes a new job manager based on the provided configuration.
// It returns an instance of IJobManager or an error if the backend is invalid.
func NewJobManager(config *Config) (jobManager IJobManager, err error) {
	switch config.Backend {
	case "redis":
		return NewRedisJobManager(config)
	default:
		return nil, ErrInvalidBackend
	}
}
