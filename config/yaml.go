package config

import (
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type yamlDecoder struct {
	cfg *Config
}

// LoadFromYAML loads configuration data from a YAML byte slice into a Configure object using a Config.
func LoadFromYAML(cfg *Config, config Configure, data []byte) (err error) {
	var node yaml.Node
	decoder := yamlDecoder{cfg: cfg}

	if err = yaml.Unmarshal(data, &node); err != nil {
		return
	}

	if cfg.YAMLFileLoading {
		if err = decoder.processYAMLNode(&node); err != nil {
			return
		}
	}

	return node.Decode(config)
}

// processYAMLNode recursively processes a YAML node for file loading and nested content handling based on the context.
func (r *yamlDecoder) processYAMLNode(node *yaml.Node) (err error) {
	if r.cfg.YAMLFileLoading {
		err = r.processYAMLFileLoading(node)
		if err != nil {
			return
		}
	}

	for _, child := range node.Content {
		err = r.processYAMLNode(child)
		if err != nil {
			return err
		}
	}

	return
}

// processYAMLFileLoading processes a YAML node, replaces scalar file references with file contents if tagged as `!file`.
func (*yamlDecoder) processYAMLFileLoading(node *yaml.Node) error {
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
