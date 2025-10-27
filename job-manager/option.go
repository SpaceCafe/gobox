package job_manager

type Option func(Manager)

func WithConfig(cfg *Config) Option {
	return func(m Manager) {

	}

}
