package logger

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name        string
		opts        []Option
		verify      func(t *testing.T, logger *DefaultLogger)
		cleanupFunc func(t *testing.T)
	}{
		{
			name: "default logger",
			verify: func(t *testing.T, logger *DefaultLogger) {
				assert.Equal(t, InfoLevel, logger.Level())
				assert.Equal(t, PlainFormat, logger.Format())
			},
		},
		{
			name: "logger with debug level",
			opts: []Option{
				WithLevel(DebugLevel),
			},
			verify: func(t *testing.T, logger *DefaultLogger) {
				var buf bytes.Buffer
				logger.logger.SetOutput(&buf)
				logger.Debug("test message")
				output := buf.String()
				assert.Contains(t, output, "DEBUG")
				assert.Contains(t, output, "test message")
				assert.Equal(t, DebugLevel, logger.Level())
			},
		},
		{
			name: "logger with info level",
			opts: []Option{
				WithLevel(InfoLevel),
			},
			verify: func(t *testing.T, logger *DefaultLogger) {
				var buf bytes.Buffer
				logger.logger.SetOutput(&buf)
				logger.Info("test message")
				output := buf.String()
				assert.Contains(t, output, "INFO")
				assert.Contains(t, output, "test message")
				assert.Equal(t, InfoLevel, logger.Level())
			},
		},
		{
			name: "logger with warn level",
			opts: []Option{
				WithLevel(WarnLevel),
			},
			verify: func(t *testing.T, logger *DefaultLogger) {
				var buf bytes.Buffer
				logger.logger.SetOutput(&buf)
				logger.Warn("test message")
				output := buf.String()
				assert.Contains(t, output, "WARN")
				assert.Contains(t, output, "test message")
				assert.Equal(t, WarnLevel, logger.Level())
			},
		},
		{
			name: "logger with error level",
			opts: []Option{
				WithLevel(ErrorLevel),
			},
			verify: func(t *testing.T, logger *DefaultLogger) {
				var buf bytes.Buffer
				logger.logger.SetOutput(&buf)
				logger.Error("test message")
				output := buf.String()
				assert.Contains(t, output, "ERROR")
				assert.Contains(t, output, "test message")
				assert.Equal(t, ErrorLevel, logger.Level())
			},
		},
		{
			name: "logger with fatal level",
			opts: []Option{
				WithLevel(FatalLevel),
			},
			verify: func(t *testing.T, logger *DefaultLogger) {
				var exitCode int
				osExit = func(code int) {
					exitCode = code
				}
				var buf bytes.Buffer
				logger.logger.SetOutput(&buf)
				logger.Fatal("test message")
				output := buf.String()
				assert.Contains(t, output, "FATAL")
				assert.Contains(t, output, "test message")
				assert.Equal(t, FatalLevel, logger.Level())
				assert.Equal(t, 1, exitCode)
			},
		},
		{
			name: "logger with JSON format",
			opts: []Option{
				WithFormat(JSONFormat),
			},
			verify: func(t *testing.T, logger *DefaultLogger) {
				var buf bytes.Buffer
				logger.logger.SetOutput(&buf)

				logger.Info("test message")
				output := buf.String()

				assert.Equal(t, JSONFormat, logger.Format())
				assert.Contains(t, output, `"level":"info"`)
				assert.Contains(t, output, `"message":"test message"`)

				buf.Reset()
			},
		},
		{
			name: "logger with Syslog format",
			opts: []Option{
				WithFormat(SyslogFormat),
			},
			verify: func(t *testing.T, logger *DefaultLogger) {
				var buf bytes.Buffer
				logger.logger.SetOutput(&buf)
				logger.Info("test\nmessage")
				output := buf.String()
				assert.Equal(t, SyslogFormat, logger.Format())
				assert.Regexp(t, `^<134>1\s.*\s+test\\nmessage\n$`, output)
			},
		},
		{
			name: "logger with custom output",
			opts: []Option{
				WithOutput("test.log"),
			},
			verify: func(t *testing.T, logger *DefaultLogger) {
				assert.NotNil(t, logger.logger)
			},
			cleanupFunc: func(t *testing.T) {
				require.NoError(t, os.Remove("test.log"))
			},
		},
		{
			name: "logger with invalid level",
			opts: []Option{
				WithLevel(Level(999)),
			},
			verify: func(t *testing.T, logger *DefaultLogger) {
				assert.Equal(t, InfoLevel, logger.Level())
			},
		},
		{
			name: "logger with invalid format",
			opts: []Option{
				WithFormat(Format(999)),
			},
			verify: func(t *testing.T, logger *DefaultLogger) {
				assert.Equal(t, PlainFormat, logger.Format())
			},
		},
		{
			name: "logger with invalid output path",
			opts: []Option{
				WithOutput("/nonexistent/directory/log.txt"),
			},
			verify: func(t *testing.T, logger *DefaultLogger) {
				assert.NotNil(t, logger)
				assert.NotNil(t, logger.logger)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.cleanupFunc != nil {
				defer tt.cleanupFunc(t)
			}

			got := New(tt.opts...)
			if tt.verify != nil {
				tt.verify(t, got)
			}
		})
	}
}
