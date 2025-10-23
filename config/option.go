package config

import (
	logger2 "github.com/spacecafe/gobox/logger"
)

// Option is a functional option that can be applied to any logger implementation.
// Implementations should accept these during construction.
type Option func(*Config)

// WithLogger sets a custom logger in the Config configuration. This allows overriding the default logging implementation.
func WithLogger(logger logger2.Logger) Option {
	return func(c *Config) {
		c.Logger = logger
	}
}

// WithConfigFilePath sets the config file path in the Config structure using the provided filepath value.
func WithConfigFilePath(filepath string) Option {
	return func(c *Config) {
		c.ConfigFilePath = filepath
	}
}

// WithEnvPrefix sets the environment variable prefix.
func WithEnvPrefix(prefix string) Option {
	return func(c *Config) {
		c.EnvPrefix = prefix
	}
}

// WithMetadata sets the metadata field in Config to the provided *types.Metadata value.
func WithMetadata(metadata *Metadata) Option {
	return func(c *Config) {
		c.Metadata = metadata
	}
}

// WithEnvAliases configures environment variable alias mappings to use when resolving environment values.
func WithEnvAliases(envAliases map[string]string) Option {
	return func(c *Config) {
		c.EnvAliases = envAliases
	}
}

// WithFlags is an option that enables flag-based configuration.
func WithFlags() Option {
	return func(c *Config) {
		c.Flags = true
	}
}

// WithEnvFileLoading enables the loading of environment variables from a specified file.
func WithEnvFileLoading() Option {
	return func(c *Config) {
		c.EnvFileLoading = true
	}
}

// WithYAMLEnvExpansion enables YAML environment variable expansion.
func WithYAMLEnvExpansion() Option {
	return func(c *Config) {
		c.YAMLEnvExpansion = true
	}
}

// WithYAMLFileLoading enables YAML file loading syntax `!file` to replace the value with the file content.
func WithYAMLFileLoading() Option {
	return func(c *Config) {
		c.YAMLFileLoading = true
	}
}
