package logger

// Option is a functional option that can be applied to any logger implementation.
// Implementations should accept these during construction.
type Option func(ConfigurableLogger)

// WithLevel sets the log level.
func WithLevel(level Level) Option {
	return func(l ConfigurableLogger) {
		_ = l.SetLevel(level)
	}
}

// WithFormat sets the log format.
func WithFormat(format Format) Option {
	return func(l ConfigurableLogger) {
		_ = l.SetFormat(format)
	}
}

// WithFileOutput sets the output destination.
func WithFileOutput(filename string) Option {
	return func(l ConfigurableLogger) {
		_ = l.SetFileOutput(filename)
	}
}

// WithConfig applies a full configuration.
func WithConfig(cfg *Config) Option {
	return func(l ConfigurableLogger) {
		_ = l.SetLevel(cfg.Level)
		_ = l.SetFormat(cfg.Format)
		_ = l.SetFileOutput(cfg.Output)
	}
}
