package config

import (
	"flag"
	"fmt"
)

// setupFlags sets up command-line flags and defines their usage messages based on the application's metadata.
func setupFlags(cfg *Config) {
	flag.StringVar(&cfg.ConfigFilePath, "config", "", "path to configuration file")
	flag.Usage = func() {
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n\n", cfg.Metadata.AppName)
		hasMetadata := false
		if cfg.Metadata.Description != "" {
			_, _ = fmt.Fprintf(flag.CommandLine.Output(), "Description: %s\n", cfg.Metadata.Description)
			hasMetadata = true
		}
		if cfg.Metadata.Version != "" {
			_, _ = fmt.Fprintf(flag.CommandLine.Output(), "Version: %s\n", cfg.Metadata.Version)
			hasMetadata = true
		}
		if cfg.Metadata.Author != "" {
			_, _ = fmt.Fprintf(flag.CommandLine.Output(), "Author: %s\n", cfg.Metadata.Author)
			hasMetadata = true
		}
		if cfg.Metadata.License != "" {
			_, _ = fmt.Fprintf(flag.CommandLine.Output(), "License: %s\n", cfg.Metadata.License)
			hasMetadata = true
		}
		if cfg.Metadata.URL != "" {
			_, _ = fmt.Fprintf(flag.CommandLine.Output(), "URL: %s\n", cfg.Metadata.URL)
			hasMetadata = true
		}
		if hasMetadata {
			_, _ = flag.CommandLine.Output().Write([]byte("\n"))
		}
		flag.PrintDefaults()
	}
	flag.Parse()
}
