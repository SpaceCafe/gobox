package types

// WithConfigFilePath represents a type for specifying the path to a configuration file as a string.
type WithConfigFilePath string

// WithEnvAliases is a map type that associates YAML path names with their respective alias names.
type WithEnvAliases map[string]string

// WithEnvFileLoading enables the functionality for loading data from files in environment variables via `EXAMPLE_FILE`.
type WithEnvFileLoading struct{}

// WithEnvPrefix defines a type for specifying an environment variable prefix for configurations.
type WithEnvPrefix string

// WithFlags enables the functionality for setting up and parsing flags.
type WithFlags struct{}

// WithYAMLEnvExpansion enables the environment variable expansion in configuration processing.
type WithYAMLEnvExpansion struct{}

// WithYAMLFileLoading enables the functionality for loading data from files in YAML via `!file /path/to/file`.
type WithYAMLFileLoading struct{}
