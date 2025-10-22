package types

import (
	"errors"
)

var (
	ErrInvalidOutput = errors.New("invalid output format")
)

type Config struct {
	Level  Level  `json:"level" yaml:"level" mapstructure:"level"`
	Format Format `json:"format" yaml:"format" mapstructure:"format"`
	Output string `json:"output" yaml:"output" mapstructure:"output"`
}

func (r *Config) SetDefaults() {
	r.Level = DebugLevel
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
