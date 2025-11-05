package logger_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/spacecafe/gobox/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		opts        []logger.Option
		verify      func(t *testing.T, log *logger.DefaultLogger)
		cleanupFunc func(t *testing.T)
	}{
		{
			name: "default logger",
			verify: func(t *testing.T, log *logger.DefaultLogger) {
				t.Helper()

				assert.Equal(t, logger.InfoLevel, log.Level())
				assert.Equal(t, logger.PlainFormat, log.Format())
			},
		},
		{
			name: "logger with debug level",
			opts: []logger.Option{
				logger.WithLevel(logger.DebugLevel),
			},
			verify: func(t *testing.T, log *logger.DefaultLogger) {
				t.Helper()

				var buf bytes.Buffer
				log.SetOutput(&buf)
				log.Debug("test message")

				output := buf.String()
				assert.Contains(t, output, "DEBUG")
				assert.Contains(t, output, "test message")
				assert.Equal(t, logger.DebugLevel, log.Level())
			},
		},
		{
			name: "logger with info level",
			opts: []logger.Option{
				logger.WithLevel(logger.InfoLevel),
			},
			verify: func(t *testing.T, log *logger.DefaultLogger) {
				t.Helper()

				var buf bytes.Buffer
				log.SetOutput(&buf)
				log.Info("test message")

				output := buf.String()
				assert.Contains(t, output, "INFO")
				assert.Contains(t, output, "test message")
				assert.Equal(t, logger.InfoLevel, log.Level())
			},
		},
		{
			name: "logger with warn level",
			opts: []logger.Option{
				logger.WithLevel(logger.WarnLevel),
			},
			verify: func(t *testing.T, log *logger.DefaultLogger) {
				t.Helper()

				var buf bytes.Buffer
				log.SetOutput(&buf)
				log.Warn("test message")

				output := buf.String()
				assert.Contains(t, output, "WARN")
				assert.Contains(t, output, "test message")
				assert.Equal(t, logger.WarnLevel, log.Level())
			},
		},
		{
			name: "logger with error level",
			opts: []logger.Option{
				logger.WithLevel(logger.ErrorLevel),
			},
			verify: func(t *testing.T, log *logger.DefaultLogger) {
				t.Helper()

				var buf bytes.Buffer
				log.SetOutput(&buf)
				log.Error("test message")

				output := buf.String()
				assert.Contains(t, output, "ERROR")
				assert.Contains(t, output, "test message")
				assert.Equal(t, logger.ErrorLevel, log.Level())
			},
		},
		{
			name: "logger with fatal level",
			opts: []logger.Option{
				logger.WithLevel(logger.FatalLevel),
			},
			verify: func(t *testing.T, log *logger.DefaultLogger) {
				t.Helper()

				var exitCode int

				logger.OsExit = func(code int) {
					exitCode = code
				}

				var buf bytes.Buffer
				log.SetOutput(&buf)
				log.Fatal("test message") //nolint:revive // OsExit is mocked.

				output := buf.String()
				assert.Contains(t, output, "FATAL")
				assert.Contains(t, output, "test message")
				assert.Equal(t, logger.FatalLevel, log.Level())
				assert.Equal(t, 1, exitCode)
			},
		},
		{
			name: "logger with JSON format",
			opts: []logger.Option{
				logger.WithFormat(logger.JSONFormat),
			},
			verify: func(t *testing.T, log *logger.DefaultLogger) {
				t.Helper()

				var buf bytes.Buffer
				log.SetOutput(&buf)
				log.Info("test message")

				output := buf.String()

				assert.Equal(t, logger.JSONFormat, log.Format())
				assert.Contains(t, output, `"level":"info"`)
				assert.Contains(t, output, `"message":"test message"`)

				buf.Reset()
			},
		},
		{
			name: "logger with Syslog format",
			opts: []logger.Option{
				logger.WithFormat(logger.SyslogFormat),
			},
			verify: func(t *testing.T, log *logger.DefaultLogger) {
				t.Helper()

				var buf bytes.Buffer
				log.SetOutput(&buf)
				log.Info("test\nmessage")

				output := buf.String()

				assert.Equal(t, logger.SyslogFormat, log.Format())
				assert.Regexp(t, `^<134>1\s.*\s+test\\nmessage\n$`, output)
			},
		},
		{
			name: "logger with custom output",
			opts: []logger.Option{
				logger.WithFileOutput("test.log"),
			},
			verify: func(t *testing.T, _ *logger.DefaultLogger) {
				t.Helper()

				assert.FileExists(t, "test.log")
			},
			cleanupFunc: func(t *testing.T) {
				t.Helper()

				require.NoError(t, os.Remove("test.log"))
			},
		},
		{
			name: "logger with invalid level",
			opts: []logger.Option{
				logger.WithLevel(logger.Level(999)),
			},
			verify: func(t *testing.T, log *logger.DefaultLogger) {
				t.Helper()

				assert.Equal(t, logger.InfoLevel, log.Level())
			},
		},
		{
			name: "logger with invalid format",
			opts: []logger.Option{
				logger.WithFormat(logger.Format(999)),
			},
			verify: func(t *testing.T, log *logger.DefaultLogger) {
				t.Helper()

				assert.Equal(t, logger.PlainFormat, log.Format())
			},
		},
		{
			name: "logger with invalid output path",
			opts: []logger.Option{
				logger.WithFileOutput("/nonexistent/directory/log.txt"),
			},
			verify: func(t *testing.T, log *logger.DefaultLogger) {
				t.Helper()

				assert.NotNil(t, log)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.cleanupFunc != nil {
				defer tt.cleanupFunc(t)
			}

			got := logger.New(tt.opts...)
			if tt.verify != nil {
				tt.verify(t, got)
			}
		})
	}
}
