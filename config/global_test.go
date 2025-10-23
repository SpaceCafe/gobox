package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockConfig struct {
	Host   string `yaml:"host"`
	Port   int    `yaml:"port"`
	Toggle bool   `yaml:"toggle"`
	Name   struct {
		First string `yaml:"first"`
		Last  string `yaml:"last"`
	} `yaml:"name"`
	Tags []string `yaml:"tags"`
}

func (r *mockConfig) SetDefaults() {
	r.Host = "127.0.0.1"
	r.Port = 8080
}

func (r *mockConfig) Validate() error { return nil }

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name        string
		config      Configure
		opts        []Option
		wantErr     bool
		verify      func(t *testing.T, config Configure)
		setupFunc   func(t *testing.T) string
		cleanupFunc func(t *testing.T)
	}{
		{
			name:    "valid config",
			config:  &mockConfig{},
			wantErr: false,
			verify: func(t *testing.T, config Configure) {
				tc := config.(*mockConfig)
				assert.Equal(t, "yaml.example.com", tc.Host)
				assert.Equal(t, 7000, tc.Port)
				assert.True(t, tc.Toggle)
			},
			setupFunc: func(t *testing.T) string {
				tmpDir := t.TempDir()
				yamlContent := `host: yaml.example.com
port: 7000
toggle: true
`
				configFile := filepath.Join(tmpDir, "config.yaml")
				err := os.WriteFile(configFile, []byte(yamlContent), 0o600)
				require.NoError(t, err)
				return configFile
			},
		},
		{
			name:    "config file not found",
			config:  &mockConfig{},
			wantErr: true,
			setupFunc: func(_ *testing.T) string {
				return "/non/existent/path/config.yaml"
			},
		},
		{
			name:    "invalid yaml format",
			config:  &mockConfig{},
			wantErr: true,
			setupFunc: func(t *testing.T) string {
				tmpDir := t.TempDir()
				invalidYaml := `host: yaml.example.com
port: invalid_port
toggle: [invalid
`
				configFile := filepath.Join(tmpDir, "config.yaml")
				err := os.WriteFile(configFile, []byte(invalidYaml), 0o600)
				require.NoError(t, err)
				return configFile
			},
		},
		{
			name:    "empty config file",
			config:  &mockConfig{},
			wantErr: false,
			verify: func(t *testing.T, config Configure) {
				tc := config.(*mockConfig)
				assert.Equal(t, "127.0.0.1", tc.Host)
				assert.Equal(t, 8080, tc.Port)
				assert.False(t, tc.Toggle)
			},
			setupFunc: func(t *testing.T) string {
				tmpDir := t.TempDir()
				configFile := filepath.Join(tmpDir, "config.yaml")
				err := os.WriteFile(configFile, []byte(""), 0o600)
				require.NoError(t, err)
				return configFile
			},
		},
		{
			name:    "partial config uses defaults",
			config:  &mockConfig{},
			wantErr: false,
			verify: func(t *testing.T, config Configure) {
				tc := config.(*mockConfig)
				assert.Equal(t, "partial.example.com", tc.Host)
				assert.Equal(t, 8080, tc.Port)
				assert.False(t, tc.Toggle)
			},
			setupFunc: func(t *testing.T) string {
				tmpDir := t.TempDir()
				yamlContent := `host: partial.example.com
`
				configFile := filepath.Join(tmpDir, "config.yaml")
				err := os.WriteFile(configFile, []byte(yamlContent), 0o600)
				require.NoError(t, err)
				return configFile
			},
		},
		{
			name:    "config with extra fields",
			config:  &mockConfig{},
			wantErr: false,
			verify: func(t *testing.T, config Configure) {
				tc := config.(*mockConfig)
				assert.Equal(t, "extra.example.com", tc.Host)
				assert.Equal(t, 9000, tc.Port)
			},
			setupFunc: func(t *testing.T) string {
				tmpDir := t.TempDir()
				yamlContent := `host: extra.example.com
port: 9000
toggle: false
extra_field: this_should_be_ignored
`
				configFile := filepath.Join(tmpDir, "config.yaml")
				err := os.WriteFile(configFile, []byte(yamlContent), 0o600)
				require.NoError(t, err)
				return configFile
			},
		},
		{
			name:    "no config flag provided",
			config:  &mockConfig{},
			wantErr: false,
			verify: func(t *testing.T, config Configure) {
				tc := config.(*mockConfig)
				// Should use defaults when no config file is specified
				assert.Equal(t, "127.0.0.1", tc.Host)
				assert.Equal(t, 8080, tc.Port)
			},
			setupFunc: func(_ *testing.T) string {
				return "" // No config file
			},
		},
		{
			name:    "config with zero values",
			config:  &mockConfig{},
			wantErr: false,
			verify: func(t *testing.T, config Configure) {
				tc := config.(*mockConfig)
				assert.Equal(t, "", tc.Host)
				assert.Equal(t, 0, tc.Port)
				assert.False(t, tc.Toggle)
			},
			setupFunc: func(t *testing.T) string {
				tmpDir := t.TempDir()
				yamlContent := `host: ""
port: 0
toggle: false
`
				configFile := filepath.Join(tmpDir, "config.yaml")
				err := os.WriteFile(configFile, []byte(yamlContent), 0o600)
				require.NoError(t, err)
				return configFile
			},
		},
		{
			name:    "config with special characters",
			config:  &mockConfig{},
			wantErr: false,
			verify: func(t *testing.T, config Configure) {
				tc := config.(*mockConfig)
				assert.Equal(t, "special-chars_123.example.com", tc.Host)
				assert.Equal(t, 8888, tc.Port)
			},
			setupFunc: func(t *testing.T) string {
				tmpDir := t.TempDir()
				yamlContent := `host: special-chars_123.example.com
port: 8888
toggle: true
`
				configFile := filepath.Join(tmpDir, "config.yaml")
				err := os.WriteFile(configFile, []byte(yamlContent), 0o600)
				require.NoError(t, err)
				return configFile
			},
		},
		{
			name:    "config file with comments",
			config:  &mockConfig{},
			wantErr: false,
			verify: func(t *testing.T, config Configure) {
				tc := config.(*mockConfig)
				assert.Equal(t, "commented.example.com", tc.Host)
				assert.Equal(t, 5555, tc.Port)
				assert.True(t, tc.Toggle)
			},
			setupFunc: func(t *testing.T) string {
				tmpDir := t.TempDir()
				yamlContent := `# This is a comment
host: commented.example.com
# Another comment
port: 5555
toggle: true # inline comment
`
				configFile := filepath.Join(tmpDir, "config.yaml")
				err := os.WriteFile(configFile, []byte(yamlContent), 0o600)
				require.NoError(t, err)
				return configFile
			},
		},
		{
			name:    "environment variable",
			config:  &mockConfig{},
			wantErr: false,
			verify: func(t *testing.T, config Configure) {
				tc := config.(*mockConfig)
				assert.Equal(t, "env.example.com", tc.Host)
				assert.Equal(t, 5555, tc.Port)
				assert.True(t, tc.Toggle)
			},
			setupFunc: func(_ *testing.T) string {
				_ = os.Setenv("HOST", "env.example.com")
				_ = os.Setenv("PORT", "5555")
				_ = os.Setenv("TOGGLE", "true")
				return ""
			},
			cleanupFunc: func(_ *testing.T) {
				_ = os.Unsetenv("HOST")
				_ = os.Unsetenv("PORT")
				_ = os.Unsetenv("TOGGLE")
			},
		},
		{
			name:    "config file with file reference",
			config:  &mockConfig{},
			opts:    []Option{WithYAMLFileLoading()},
			wantErr: false,
			verify: func(t *testing.T, config Configure) {
				tc := config.(*mockConfig)
				assert.Equal(t, "ref.example.com", tc.Host)
				assert.Equal(t, 5555, tc.Port)
				assert.False(t, tc.Toggle)
			},
			setupFunc: func(t *testing.T) string {
				tmpDir := t.TempDir()
				refFile := filepath.Join(tmpDir, "testdata.txt")
				yamlContent := `host: !file ` + refFile + `
port: 5555
`
				configFile := filepath.Join(tmpDir, "config.yaml")
				err := os.WriteFile(configFile, []byte(yamlContent), 0o600)
				require.NoError(t, err)
				err = os.WriteFile(refFile, []byte("ref.example.com"), 0o600)
				require.NoError(t, err)
				return configFile
			},
		},
		{
			name:    "config file with env expansion",
			config:  &mockConfig{},
			opts:    []Option{WithYAMLEnvExpansion()},
			wantErr: false,
			verify: func(t *testing.T, config Configure) {
				tc := config.(*mockConfig)
				assert.Equal(t, "env.example.com", tc.Host)
				assert.Equal(t, 5555, tc.Port)
				assert.False(t, tc.Toggle)
			},
			setupFunc: func(t *testing.T) string {
				tmpDir := t.TempDir()
				yamlContent := `host: ${HOST}
port: $PORT
`
				configFile := filepath.Join(tmpDir, "config.yaml")
				err := os.WriteFile(configFile, []byte(yamlContent), 0o600)
				require.NoError(t, err)
				_ = os.Setenv("HOST", "env.example.com")
				_ = os.Setenv("PORT", "5555")
				return configFile
			},
			cleanupFunc: func(_ *testing.T) {
				_ = os.Unsetenv("HOST")
				_ = os.Unsetenv("PORT")
			},
		},
		{
			name:    "config file and env",
			config:  &mockConfig{},
			opts:    []Option{WithEnvPrefix("APP")},
			wantErr: false,
			verify: func(t *testing.T, config Configure) {
				tc := config.(*mockConfig)
				assert.Equal(t, "env.example.com", tc.Host)
				assert.Equal(t, 5555, tc.Port)
				assert.False(t, tc.Toggle)
			},
			setupFunc: func(t *testing.T) string {
				tmpDir := t.TempDir()
				yamlContent := `host: yaml.example.com
port: 5555
`
				configFile := filepath.Join(tmpDir, "config.yaml")
				err := os.WriteFile(configFile, []byte(yamlContent), 0o600)
				require.NoError(t, err)
				_ = os.Setenv("APP_HOST", "env.example.com")
				return configFile
			},
			cleanupFunc: func(_ *testing.T) {
				_ = os.Unsetenv("APP_HOST")
			},
		},
		{
			name:    "env with prefix and nested",
			config:  &mockConfig{},
			opts:    []Option{WithEnvPrefix("APP")},
			wantErr: false,
			verify: func(t *testing.T, config Configure) {
				tc := config.(*mockConfig)
				assert.Equal(t, "Jane", tc.Name.First)
				assert.Equal(t, "Doe", tc.Name.Last)
			},
			setupFunc: func(_ *testing.T) string {
				_ = os.Setenv("APP_NAME_FIRST", "Jane")
				_ = os.Setenv("APP_NAME_LAST", "Doe")
				return ""
			},
			cleanupFunc: func(_ *testing.T) {
				_ = os.Unsetenv("APP_NAME_FIRST")
				_ = os.Unsetenv("APP_NAME_LAST")
			},
		},
		{
			name:    "env tags",
			config:  &mockConfig{},
			wantErr: false,
			verify: func(t *testing.T, config Configure) {
				tc := config.(*mockConfig)
				assert.Equal(t, []string{"foo", "bar", "baz"}, tc.Tags)
			},
			setupFunc: func(_ *testing.T) string {
				_ = os.Setenv("TAGS", "foo,bar,baz")
				return ""
			},
			cleanupFunc: func(_ *testing.T) {
				_ = os.Unsetenv("TAGS")
			},
		},
		{
			name:    "env aliases",
			config:  &mockConfig{},
			opts:    []Option{WithEnvAliases(Aliases{"name.first": "FIRSTNAME"})},
			wantErr: false,
			verify: func(t *testing.T, config Configure) {
				tc := config.(*mockConfig)
				assert.Equal(t, "John", tc.Name.First)
			},
			setupFunc: func(_ *testing.T) string {
				_ = os.Setenv("FIRSTNAME", "John")
				return ""
			},
			cleanupFunc: func(_ *testing.T) {
				_ = os.Unsetenv("FIRSTNAME")
			},
		},
		{
			name:    "env file loading via alias",
			config:  &mockConfig{},
			opts:    []Option{WithEnvAliases(Aliases{"name.first": "FIRSTNAME"}), WithEnvFileLoading()},
			wantErr: false,
			verify: func(t *testing.T, config Configure) {
				tc := config.(*mockConfig)
				assert.Equal(t, "John", tc.Name.First)
			},
			setupFunc: func(_ *testing.T) string {
				tmpDir := t.TempDir()
				refFile := filepath.Join(tmpDir, "test.txt")
				err := os.WriteFile(refFile, []byte("John \n"), 0o600)
				require.NoError(t, err)
				_ = os.Setenv("FIRSTNAME_FILE", refFile)
				return ""
			},
			cleanupFunc: func(_ *testing.T) {
				_ = os.Unsetenv("FIRSTNAME_FILE")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupFunc != nil {
				tt.opts = append(tt.opts, WithConfigFilePath(tt.setupFunc(t)))
			}

			if tt.cleanupFunc != nil {
				defer tt.cleanupFunc(t)
			}

			gotErr := LoadConfig(tt.config, tt.opts...)
			assert.Equal(t, tt.wantErr, gotErr != nil, gotErr)
			if !tt.wantErr && tt.verify != nil {
				tt.verify(t, tt.config)
			}
		})
	}
}
