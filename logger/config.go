package logger

import (
	"errors"

	"github.com/spacecafe/gobox/config"
)

var (
	_ config.Configure = (*Config)(nil)

	ErrInvalidOutput = errors.New("invalid output format")
)

type Config struct {
	Level  Level  `json:"level" yaml:"level" mapstructure:"level"`
	Format Format `json:"format" yaml:"format" mapstructure:"format"`
	Output string `json:"output" yaml:"output" mapstructure:"output"`
}

func (r *Config) SetDefaults() {
	r.Level = InfoLevel
	r.Format = PlainFormat
	r.Output = "/dev/stderr"
}

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
