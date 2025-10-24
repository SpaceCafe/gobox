package logger

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
)

// DefaultLogger is the default implementation of a logger with configurable format, level, and message output.
type DefaultLogger struct {
	appName string
	logger  *log.Logger
	format  Format
	level   Level
	output  func(*Entry, int)
}

// Ensure DefaultLogger implements the interface.
var _ AdvancedLogger = (*DefaultLogger)(nil)

// New creates a new default DefaultLogger instance.
func New(opts ...Option) *DefaultLogger {
	cfg := &Config{}
	cfg.SetDefaults()

	l := &DefaultLogger{
		appName: filepath.Base(os.Args[0]),
		logger:  log.New(os.Stderr, "", 0),
		level:   cfg.Level,
	}

	_ = l.SetFormat(cfg.Format)

	for _, opt := range opts {
		opt(l)
	}

	return l
}

// Debug writes debug level messages.
func (r *DefaultLogger) Debug(v ...any) {
	//nolint:mnd // Skips over this method in stack trace.
	r.Output(DebugLevel, 2, nil, v...)
}

// Debugf writes debug level messages using formatted string.
func (r *DefaultLogger) Debugf(format string, v ...any) {
	//nolint:mnd // Skips over this method in stack trace.
	r.Output(DebugLevel, 2, &format, v...)
}

// Info writes info level messages.
func (r *DefaultLogger) Info(v ...any) {
	//nolint:mnd // Skips over this method in stack trace.
	r.Output(InfoLevel, 2, nil, v...)
}

// Infof writes info level messages using formatted string.
func (r *DefaultLogger) Infof(format string, v ...any) {
	//nolint:mnd // Skips over this method in stack trace.
	r.Output(InfoLevel, 2, &format, v...)
}

// Warn writes warning level messages.
func (r *DefaultLogger) Warn(v ...any) {
	//nolint:mnd // Skips over this method in stack trace.
	r.Output(WarnLevel, 2, nil, v...)
}

// Warnf writes warning level messages using formatted string.
func (r *DefaultLogger) Warnf(format string, v ...any) {
	//nolint:mnd // Skips over this method in stack trace.
	r.Output(WarnLevel, 2, &format, v...)
}

// Warning writes warning level messages as an alias for the Warn method.
func (r *DefaultLogger) Warning(v ...any) {
	//nolint:mnd // Skips over this method in stack trace.
	r.Output(WarnLevel, 2, nil, v...)
}

// Warningf writes warning level messages using a formatted string. It acts as an alias for the Warnf method.
func (r *DefaultLogger) Warningf(format string, v ...any) {
	//nolint:mnd // Skips over this method in stack trace.
	r.Output(WarnLevel, 2, &format, v...)
}

// Error writes error level messages.
func (r *DefaultLogger) Error(v ...any) {
	//nolint:mnd // Skips over this method in stack trace.
	r.Output(ErrorLevel, 2, nil, v...)
}

// Errorf writes error level messages using formatted string.
func (r *DefaultLogger) Errorf(format string, v ...any) {
	//nolint:mnd // Skips over this method in stack trace.
	r.Output(ErrorLevel, 2, &format, v...)
}

// Fatal writes fatal level messages and terminates the application.
func (r *DefaultLogger) Fatal(v ...any) {
	//nolint:mnd // Skips over this method in stack trace.
	r.Output(FatalLevel, 2, nil, v...)
	osExit(1)
}

// Fatalf writes fatal level messages using formatted string and terminates the application.
func (r *DefaultLogger) Fatalf(format string, v ...any) {
	//nolint:mnd // Skips over this method in stack trace.
	r.Output(FatalLevel, 2, &format, v...)
	osExit(1)
}

// Level returns the current logging level of the logger.
func (r *DefaultLogger) Level() Level {
	return r.level
}

// SetLevel changes the logging level of the logger to the specified level.
func (r *DefaultLogger) SetLevel(level Level) error {
	if level < DebugLevel || level > FatalLevel {
		return ErrInvalidLevel
	}
	r.level = level
	return nil
}

// Format returns the current output format setting of the logger.
func (r *DefaultLogger) Format() Format {
	return r.format
}

// SetFormat changes the output format of the logger to the specified format.
func (r *DefaultLogger) SetFormat(format Format) (err error) {
	switch format {
	case PlainFormat:
		r.output = r.outputPlain
	case JSONFormat:
		r.output = r.outputJSON
	case SyslogFormat:
		r.output = r.outputSyslog
	default:
		return ErrInvalidFormat
	}
	r.format = format
	return nil
}

// SetOutput opens or create the given file and set it as the new logging destination.
func (r *DefaultLogger) SetOutput(filename string) error {
	// Using 0o600 permission to ensure only the owner can read and write the log file for security
	//nolint:mnd // Permission flag is a valid constant
	file, err := os.OpenFile(path.Clean(filename), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if r.logger != nil {
		r.logger.SetOutput(file)
	}
	return err
}

// Output is used to log messages with a specific level and format.
func (r *DefaultLogger) Output(level Level, calldepth int, format *string, v ...any) {
	// If the current logger's level is higher than the given level, return immediately without logging anything.
	if r.level > level {
		return
	}

	var (
		file string
		line int
		ok   bool
		l    *Entry
	)

	if r.level == DebugLevel {
		// Get information about the location of the logging call using runtime.Caller function with the provided calldepth.
		// If it fails, set default values for filename and line number.
		_, file, line, ok = runtime.Caller(calldepth)
		if !ok {
			file = "???"
			line = 0
		}

		// Shorten the filename to only include the last part after the final '/'.
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
	}

	// Create a new Entry with the given level, filename, line number, and message.
	// If a format is nil, use fmt.Sprint to create the log message; otherwise, use fmt.Sprintf with the provided format string.
	if format == nil {
		if len(v) == 1 {
			l = NewEntry(level, file, line, v[0])
		} else {
			l = NewEntry(level, file, line, fmt.Sprint(v...))
		}
	} else {
		l = NewEntry(level, file, line, fmt.Sprintf(*format, v...))
	}

	// Pass the Entry to the output function for further processing and logging.
	calldepth++
	r.output(l, calldepth)
}
