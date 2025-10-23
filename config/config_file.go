package config

import (
	"os"
	"path/filepath"
)

// lookupConfigFile searches for a configuration file in common directories and updates Config with its path.
// It checks user configuration, system directories, and local directories using predefined file name patterns.
func (r *Config) lookupConfigFile() error {
	paths := []string{"."}
	if dir, err := os.UserConfigDir(); err == nil {
		paths = append(paths, filepath.Join(dir, r.Metadata.Slug))
	}
	paths = append(paths, filepath.Join(SysConfigDir(), r.Metadata.Slug))

	fileNames := []string{
		r.Metadata.Slug + ".yml",
		r.Metadata.Slug + ".yaml",
		"config.yml",
		"config.yaml",
	}

	for _, path := range paths {
		for _, name := range fileNames {
			filePath := filepath.Join(path, name)
			if _, err := os.Stat(filePath); err == nil {
				r.ConfigFilePath = filePath
				return nil
			}
		}
	}

	return ErrConfigFileNotFound
}

// loadConfigFile reads the configuration file specified in the Config and returns its content as a byte slice.
// Expands environment variables in YAML content if WithYAMLEnvExpansion is enabled in the Config.
func (r *Config) loadConfigFile() ([]byte, error) {
	if r.ConfigFilePath == "" {
		if err := r.lookupConfigFile(); err != nil {
			return []byte{}, err
		}
	}

	data, err := os.ReadFile(r.ConfigFilePath)
	if err != nil {
		return []byte{}, err
	}

	if r.YAMLEnvExpansion {
		data = []byte(os.ExpandEnv(string(data)))
	}

	return data, nil
}
