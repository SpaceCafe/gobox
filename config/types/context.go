package types

import (
	"os"
	"path/filepath"

	"github.com/spacecafe/gobox/logger"
	loggertypes "github.com/spacecafe/gobox/logger/types"
)

// ConfigContext represents the configuration context for the config package.
type ConfigContext struct {
	Log                  loggertypes.Logger
	Metadata             *Metadata
	EnvAliases           map[string]string
	EnvPrefix            string
	ConfigFilePath       string
	WithYAMLEnvExpansion bool
	WithYAMLFileLoading  bool
	WithEnvFileLoading   bool
}

// SetDefaults initializes the ConfigContext with default values for metadata, logger, and configuration options.
func (s *ConfigContext) SetDefaults() {
	name := filepath.Base(os.Args[0])
	if name == "" {
		name = "app"
	}

	s.Log = logger.Default()
	s.Metadata = &Metadata{
		AppName: name,
		Slug:    name,
	}
	s.EnvAliases = make(map[string]string)
	s.WithYAMLEnvExpansion = false
	s.WithYAMLFileLoading = false
	s.WithEnvFileLoading = false
}
