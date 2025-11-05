package logger

type Config struct {
	Level  Level  `json:"level"  mapstructure:"level"  yaml:"level"`
	Format Format `json:"format" mapstructure:"format" yaml:"format"`
	Output string `json:"output" mapstructure:"output" yaml:"output"`
}

// SetDefaults initializes the default values for the relevant fields in the struct.
func (r *Config) SetDefaults() {
	r.Level = InfoLevel
	r.Format = PlainFormat
	r.Output = "/dev/stderr"
}

// Validate ensures the all necessary configurations are filled and within valid confines.
func (r *Config) Validate() error {
	if r.Level < DebugLevel || r.Level > FatalLevel {
		return ErrInvalidLevel
	}

	if r.Format < PlainFormat || r.Format > SyslogFormat {
		return ErrInvalidFormat
	}

	if r.Output == "" {
		return ErrInvalidOutput
	}

	return nil
}
