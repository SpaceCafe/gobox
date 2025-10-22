package config

import (
	"flag"
	"fmt"

	"github.com/spacecafe/gobox/config/types"
)

// setupFlags sets up command-line flags and defines their usage messages based on the application's metadata.
func setupFlags(ctx *types.ConfigContext) {
	flag.StringVar(&ctx.ConfigFilePath, "config", "", "path to configuration file")
	flag.Usage = func() {
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n\n", ctx.Metadata.AppName)
		hasMetadata := false
		if ctx.Metadata.Description != "" {
			_, _ = fmt.Fprintf(flag.CommandLine.Output(), "Description: %s\n", ctx.Metadata.Description)
			hasMetadata = true
		}
		if ctx.Metadata.Version != "" {
			_, _ = fmt.Fprintf(flag.CommandLine.Output(), "Version: %s\n", ctx.Metadata.Version)
			hasMetadata = true
		}
		if ctx.Metadata.Author != "" {
			_, _ = fmt.Fprintf(flag.CommandLine.Output(), "Author: %s\n", ctx.Metadata.Author)
			hasMetadata = true
		}
		if ctx.Metadata.License != "" {
			_, _ = fmt.Fprintf(flag.CommandLine.Output(), "License: %s\n", ctx.Metadata.License)
			hasMetadata = true
		}
		if ctx.Metadata.URL != "" {
			_, _ = fmt.Fprintf(flag.CommandLine.Output(), "URL: %s\n", ctx.Metadata.URL)
			hasMetadata = true
		}
		if hasMetadata {
			_, _ = flag.CommandLine.Output().Write([]byte("\n"))
		}
		flag.PrintDefaults()
	}
	flag.Parse()
}
