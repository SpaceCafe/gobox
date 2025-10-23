package config

import (
	"os"
	"path/filepath"

	"github.com/spacecafe/gobox/logger"
)

// Config contain configuration settings, including logging, metadata, environment, YAML, and file loading preferences.
type Config struct {
	Logger           logger.Logger
	ConfigFilePath   string
	EnvPrefix        string
	Metadata         *Metadata
	EnvAliases       Aliases
	EnvFileLoading   bool
	Flags            bool
	YAMLEnvExpansion bool
	YAMLFileLoading  bool
}

type Aliases map[string]string

// SetDefaults initializes the Config with default values for metadata, logger, and configuration options.
func (r *Config) SetDefaults() {
	name := filepath.Base(os.Args[0])
	if name == "" {
		name = "app"
	}

	r.Logger = logger.Default()
	r.Metadata = &Metadata{
		AppName: name,
		Slug:    name,
	}
	r.EnvAliases = make(Aliases)
	r.Flags = false
	r.YAMLEnvExpansion = false
	r.YAMLFileLoading = false
	r.EnvFileLoading = false
}
