package logger

//nolint:gochecknoglobals // std is a global logger instance that needs to be accessible throughout the package.
var std Logger = New()

// Default returns the package-level default logger.
func Default() Logger {
	return std
}

// SetDefault replaces the package-level default logger with a new one.
func SetDefault(l Logger) {
	std = l
}

// Package-level convenience functions that delegate to the default logger.

func Debug(v ...any) {
	std.Debug(v...)
}

func Debugf(format string, v ...any) {
	std.Debugf(format, v...)
}

func Info(v ...any) {
	std.Info(v...)
}

func Infof(format string, v ...any) {
	std.Infof(format, v...)
}

func Warn(v ...any) {
	std.Warn(v...)
}

func Warnf(format string, v ...any) {
	std.Warnf(format, v...)
}

func Error(v ...any) {
	std.Error(v...)
}

func Errorf(format string, v ...any) {
	std.Errorf(format, v...)
}

func Fatal(v ...any) {
	std.Fatal(v...)
}

func Fatalf(format string, v ...any) {
	std.Fatalf(format, v...)
}
