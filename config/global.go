package config

import (
	"errors"
)

// LoadConfig initializes and loads configuration into the provided Configure object using
// a configuration file and environment variables.
func LoadConfig(config Configure, opts ...Option) error {
	if !IsStructPointer(config) {
		return ErrInvalidConfig
	}

	cfg := &Config{}
	cfg.SetDefaults()

	for _, opt := range opts {
		opt(cfg)
	}

	if v, ok := config.(MetadataProvider); ok {
		cfg.Metadata = v.Metadata()
	}

	if cfg.Flags {
		setupFlags(cfg)
	}

	config.SetDefaults()

	configFile, err := cfg.loadConfigFile()
	if err != nil && !errors.Is(err, ErrConfigFileNotFound) {
		return err
	}
	if configFile != nil {
		err = LoadFromYAML(cfg, config, configFile)
		if err != nil {
			return err
		}
	}

	err = LoadFromEnv(cfg, config)
	if err != nil {
		return err
	}

	return config.Validate()
}
