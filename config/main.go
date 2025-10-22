package config

import (
	"errors"

	"github.com/spacecafe/gobox/config/types"
	"github.com/spacecafe/gobox/config/utils"
	loggertypes "github.com/spacecafe/gobox/logger/types"
)

// LoadConfig initializes and loads configuration into the provided Configure object using a configuration file and environment variables.
func LoadConfig(config types.Configure, args ...any) error {
	if !utils.IsStructPointer(config) {
		return types.ErrInvalidConfig
	}

	ctx := &types.ConfigContext{}
	ctx.SetDefaults()

	if v, ok := config.(types.MetadataProvider); ok {
		ctx.Metadata = v.Metadata()
	}

	config.SetDefaults()
	processArgs(ctx, args...)

	configFile, err := loadConfigFile(ctx)
	if err != nil && !errors.Is(err, types.ErrConfigFileNotFound) {
		return err
	}
	if configFile != nil {
		err = LoadFromYAML(ctx, config, configFile)
		if err != nil {
			return err
		}
	}

	err = LoadFromEnv(ctx, config)
	if err != nil {
		return err
	}

	return config.Validate()
}

// processArgs processes a variable number of arguments and updates the provided ConfigContext accordingly.
func processArgs(ctx *types.ConfigContext, args ...any) {
	withFlags := false
	for _, arg := range args {
		switch typedArg := arg.(type) {
		case types.Metadata:
			ctx.Metadata = &typedArg
		case *types.Metadata:
			ctx.Metadata = typedArg
		case loggertypes.Logger:
			ctx.Log = typedArg
		case types.WithConfigFilePath:
			ctx.ConfigFilePath = string(typedArg)
		case types.WithEnvAliases:
			ctx.EnvAliases = typedArg
		case types.WithEnvFileLoading:
			ctx.WithEnvFileLoading = true
		case types.WithEnvPrefix:
			ctx.EnvPrefix = string(typedArg)
		case types.WithFlags:
			withFlags = true
		case types.WithYAMLEnvExpansion:
			ctx.WithYAMLEnvExpansion = true
		case types.WithYAMLFileLoading:
			ctx.WithYAMLFileLoading = true
		}
	}
	if withFlags {
		setupFlags(ctx)
	}
}
