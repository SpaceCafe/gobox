package logger

import (
	"io"
)

// Logger is the main logging interface that implementations must satisfy.
type Logger interface {
	Debug(v ...any)
	Debugf(format string, v ...any)
	Info(v ...any)
	Infof(format string, v ...any)
	Warn(v ...any)
	Warnf(format string, v ...any)
	Error(v ...any)
	Errorf(format string, v ...any)
	Fatal(v ...any)
	Fatalf(format string, v ...any)
}

// ConfigurableLogger extends Logger with configuration methods.
type ConfigurableLogger interface {
	Logger

	Level() Level
	SetLevel(level Level) error

	Format() Format
	SetFormat(format Format) error

	SetFileOutput(filename string) error
	SetOutput(writer io.Writer)
}

// AdvancedLogger provides low-level output control.
type AdvancedLogger interface {
	ConfigurableLogger
	Warning(v ...any)
	Warningf(format string, v ...any)
	Output(level Level, calldepth int, format *string, v ...any)
}
