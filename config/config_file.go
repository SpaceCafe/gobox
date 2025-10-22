package config

import (
	"os"
	"path/filepath"

	"github.com/spacecafe/gobox/config/types"
	"github.com/spacecafe/gobox/config/utils"
)

// lookupConfigFile searches for a configuration file in common directories and updates ConfigContext with its path.
// It checks user configuration, system directories, and local directories using predefined file name patterns.
func lookupConfigFile(state *types.ConfigContext) error {
	paths := []string{"."}
	if dir, err := os.UserConfigDir(); err == nil {
		paths = append(paths, filepath.Join(dir, state.Metadata.Slug))
	}
	paths = append(paths, filepath.Join(utils.SysConfigDir(), state.Metadata.Slug))

	fileNames := []string{
		state.Metadata.Slug + ".yml",
		state.Metadata.Slug + ".yaml",
		"config.yml",
		"config.yaml",
	}

	for _, path := range paths {
		for _, name := range fileNames {
			filePath := filepath.Join(path, name)
			if _, err := os.Stat(filePath); err == nil {
				state.ConfigFilePath = filePath
				return nil
			}
		}
	}

	return types.ErrConfigFileNotFound
}

// loadConfigFile reads the configuration file specified in the ConfigContext and returns its content as a byte slice.
// Expands environment variables in YAML content if WithYAMLEnvExpansion is enabled in the ConfigContext.
func loadConfigFile(state *types.ConfigContext) ([]byte, error) {
	if state.ConfigFilePath == "" {
		if err := lookupConfigFile(state); err != nil {
			return []byte{}, err
		}
	}

	data, err := os.ReadFile(state.ConfigFilePath)
	if err != nil {
		return []byte{}, err
	}

	if state.WithYAMLEnvExpansion {
		data = []byte(os.ExpandEnv(string(data)))
	}

	return data, nil
}
