package logger

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"strings"
)

const (
	fileMode = 0600
)

var (
	// osExit is a variable for testing purposes.
	osExit = os.Exit

	ErrInvalidFormat = errors.New("format is invalid")
	ErrInvalidLevel  = errors.New("level is invalid")
)

// Logger allows messages with different levels/priorities to be sent to stderr/stdout
type Logger struct {
	logger *log.Logger
	format Format
	level  Level
	output func(*Item, int)
	Warn   func(...any)
	Warnf  func(string, ...any)
}

// New returns a new logger at the debug level and on stderr.
func New() *Logger {
	l := &Logger{
		logger: log.New(os.Stderr, "", 0),
	}
	l.Warn = l.Warning
	l.Warnf = l.Warningf
	_ = l.SetFormat(PlainFormat)
	_ = l.SetLevel(DebugLevel)

	return l
}

// SetFormat sets the output format for the logger.
func (r *Logger) SetFormat(format Format) error {
	switch format {
	case PlainFormat:
		r.output = r.outputPlain
	case JSONFormat:
		r.output = r.outputJSON
	default:
		return ErrInvalidFormat
	}
	r.format = format

	return nil
}

// SetLevel sets the current loglevel to the desired value.
func (r *Logger) SetLevel(level Level) error {
	if level < DebugLevel || level > FatalLevel {
		return ErrInvalidLevel
	}
	r.level = level

	return nil
}

// ParseLevel sets the current loglevel to the desired value from string identifier.
func (r *Logger) ParseLevel(name string) error {
	if level, ok := StringToLevel[strings.ToLower(name)]; ok {
		r.level = level
		return nil
	}

	return ErrInvalidLevel
}

// SetOutput opens or create the given file and set it as the new logging destination.
func (r *Logger) SetOutput(filename string) error {
	file, err := os.OpenFile(path.Clean(filename), os.O_APPEND|os.O_CREATE|os.O_WRONLY, fileMode)
	if r.logger != nil {
		r.logger.SetOutput(file)
	}

	return err
}

// Output is used to log messages with a specific level and format.
func (r *Logger) Output(level Level, calldepth int, format *string, v ...any) {
	// If the current logger's level is higher than the given level, return immediately without logging anything.
	if r.level > level {
		return
	}

	var (
		file string
		line int
		ok   bool
		l    *Item
	)

	// Get information about the location of the logging call using runtime.Caller function with provided calldepth.
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

	// Create a new Item with the given level, filename, line number, and message.
	// If format is nil, use fmt.Sprint to create the log message; otherwise, use fmt.Sprintf with the provided format string.
	if format == nil {
		if len(v) == 1 {
			l = NewItem(level, file, line, v[0])
		} else {
			l = NewItem(level, file, line, fmt.Sprint(v...))
		}
	} else {
		l = NewItem(level, file, line, fmt.Sprintf(*format, v...))
	}

	// Pass the Item to the output function for further processing and logging.
	calldepth++
	r.output(l, calldepth)
}

// outputPlain is a wrapper function to create a log entry in plain text.
func (r *Logger) outputPlain(item *Item, calldepth int) {
	_ = r.logger.Output(calldepth, item.String())
}

// outputJSON is a wrapper function to create a log entry in JSON.
func (r *Logger) outputJSON(item *Item, calldepth int) {
	text, err := item.Marshal()
	if err != nil {
		_ = r.logger.Output(calldepth, err.Error())
	} else {
		_ = r.logger.Output(calldepth, string(text))
	}
}

// Debug writes debug level messages
func (r *Logger) Debug(v ...any) {
	//nolint:gomnd
	r.Output(DebugLevel, 2, nil, v...)
}

// Debugf writes debug level messages using formatted string
func (r *Logger) Debugf(format string, v ...any) {
	//nolint:gomnd
	r.Output(DebugLevel, 2, &format, v...)
}

// Info writes info level messages
func (r *Logger) Info(v ...any) {
	//nolint:gomnd
	r.Output(InfoLevel, 2, nil, v...)
}

// Infof writes info level messages using formatted string
func (r *Logger) Infof(format string, v ...any) {
	//nolint:gomnd
	r.Output(InfoLevel, 2, &format, v...)
}

// Warning writes warning level messages
func (r *Logger) Warning(v ...any) {
	//nolint:gomnd
	r.Output(WarningLevel, 2, nil, v...)
}

// Warningf writes warning level messages using formatted string
func (r *Logger) Warningf(format string, v ...any) {
	//nolint:gomnd
	r.Output(WarningLevel, 2, &format, v...)
}

// Error writes error level messages
func (r *Logger) Error(v ...any) {
	//nolint:gomnd
	r.Output(ErrorLevel, 2, nil, v...)
}

// Errorf writes error level messages using formatted string
func (r *Logger) Errorf(format string, v ...any) {
	//nolint:gomnd
	r.Output(ErrorLevel, 2, &format, v...)
}

// Fatal writes fatal level messages and terminates the application
func (r *Logger) Fatal(v ...any) {
	//nolint:gomnd
	r.Output(FatalLevel, 2, nil, v...)
	osExit(1)
}

// Fatalf writes fatal level messages using formatted string and terminates the application
func (r *Logger) Fatalf(format string, v ...any) {
	//nolint:gomnd
	r.Output(FatalLevel, 2, &format, v...)
	osExit(1)
}
