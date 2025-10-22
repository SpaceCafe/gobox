package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spacecafe/gobox/config/types"
	"gopkg.in/yaml.v3"
)

// LoadFromYAML loads configuration data from a YAML byte slice into a Configure object using a ConfigContext.
func LoadFromYAML(ctx *types.ConfigContext, config types.Configure, data []byte) (err error) {
	var node yaml.Node
	if err = yaml.Unmarshal(data, &node); err != nil {
		return
	}

	if ctx.WithYAMLFileLoading {
		if err = processYAMLNode(ctx, &node); err != nil {
			return
		}
	}

	return node.Decode(config)
}

// processYAMLNode recursively processes a YAML node for file loading and nested content handling based on the context.
func processYAMLNode(ctx *types.ConfigContext, node *yaml.Node) (err error) {
	if ctx.WithYAMLFileLoading {
		err = processYAMLFileLoading(node)
		if err != nil {
			return
		}
	}

	for _, child := range node.Content {
		err = processYAMLNode(ctx, child)
		if err != nil {
			return err
		}
	}

	return
}

// processYAMLFileLoading processes a YAML node, replaces scalar file references with file contents if tagged as `!file`.
func processYAMLFileLoading(node *yaml.Node) error {
	if node.Kind == yaml.ScalarNode && node.Tag == "!file" {
		filePath := strings.TrimSpace(node.Value)
		content, err := os.ReadFile(filepath.Clean(filePath))
		if err != nil {
			return err
		}
		node.Value = string(content)
	}
	return nil
}
