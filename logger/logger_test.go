package logger

import (
	"bytes"
	"os"
	"testing"

	"github.com/spacecafe/gobox/logger/types"
	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {
	type args struct {
		level  types.Level
		format types.Format
	}
	type wants struct {
		errLevel  error
		errFormat error
	}
	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{"Debug", args{types.DebugLevel, types.PlainFormat}, wants{}},
		{"Info", args{types.InfoLevel, types.JSONFormat}, wants{}},
		{"Warning", args{types.WarningLevel, types.PlainFormat}, wants{}},
		{"Error", args{types.ErrorLevel, types.JSONFormat}, wants{}},
		{"Fatal", args{types.FatalLevel, types.PlainFormat}, wants{}},
		{"invalid level", args{99, types.PlainFormat}, wants{errLevel: types.ErrInvalidLevel}},
		{"invalid format", args{types.FatalLevel, 99}, wants{errFormat: types.ErrInvalidFormat}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			l := New()

			err := l.SetLevel(tt.args.level)
			assert.Equal(t, tt.wants.errLevel, err)
			if err != nil {
				return
			}

			err = l.SetFormat(tt.args.format)
			assert.Equal(t, tt.wants.errFormat, err)
			if err != nil {
				return
			}

			l.logger.SetOutput(&buf)

			// Test Debug/Debugf
			l.Debug("test message")
			if tt.args.level <= types.DebugLevel {
				assert.NotEmpty(t, buf.String())
			} else {
				assert.Empty(t, buf.String())
			}
			buf.Reset()

			l.Debugf("%s", "test message")
			if tt.args.level <= types.DebugLevel {
				assert.NotEmpty(t, buf.String())
			} else {
				assert.Empty(t, buf.String())
			}
			buf.Reset()

			// Test Info/Infof
			l.Info("test message")
			if tt.args.level <= types.InfoLevel {
				assert.NotEmpty(t, buf.String())
			} else {
				assert.Empty(t, buf.String())
			}
			buf.Reset()

			l.Infof("%s", "test message")
			if tt.args.level <= types.InfoLevel {
				assert.NotEmpty(t, buf.String())
			} else {
				assert.Empty(t, buf.String())
			}
			buf.Reset()

			// Test Warning/Warningf
			l.Warning("test message")
			if tt.args.level <= types.WarningLevel {
				assert.NotEmpty(t, buf.String())
			} else {
				assert.Empty(t, buf.String())
			}
			buf.Reset()

			l.Warningf("%s", "test message")
			if tt.args.level <= types.WarningLevel {
				assert.NotEmpty(t, buf.String())
			} else {
				assert.Empty(t, buf.String())
			}
			buf.Reset()

			// Test Error/Errorf
			l.Error("test message")
			if tt.args.level <= types.ErrorLevel {
				assert.NotEmpty(t, buf.String())
			} else {
				assert.Empty(t, buf.String())
			}
			buf.Reset()

			l.Errorf("%s", "test message")
			if tt.args.level <= types.ErrorLevel {
				assert.NotEmpty(t, buf.String())
			} else {
				assert.Empty(t, buf.String())
			}
			buf.Reset()

			// Test Fatal/Fatalf
			var exitCode int
			osExit = func(code int) {
				exitCode = code
			}
			l.Fatal("test message")
			assert.NotEmpty(t, buf.String())
			assert.NotEqual(t, 0, exitCode)
			buf.Reset()

			l.Fatalf("%s", "test message")
			assert.NotEmpty(t, buf.String())
			assert.NotEqual(t, 0, exitCode)
			buf.Reset()
		})
	}
}

func TestLogger_SetOutput(t *testing.T) {
	defer func() {
		_ = os.Truncate("testdata/test.log", 0)
		_ = os.Remove("testdata/new_test.log")
	}()

	tests := []struct {
		name     string
		filename string
	}{
		{"append file", "testdata/test.log"},
		{"create file", "testdata/new_test.log"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := New()
			err := r.SetOutput(tt.filename)
			assert.NoError(t, err)
			r.Debug("Test message")

			content, err := os.ReadFile(tt.filename)
			assert.NoError(t, err)
			assert.NotEmpty(t, content)
			assert.Contains(t, string(content), "Test message")
		})
	}
}
