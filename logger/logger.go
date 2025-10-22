package logger

import (
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/spacecafe/gobox/logger/types"
)

const (
	// plainTimeFormat sets the format using Go's reference time.
	plainTimeFormat  = "2006-01-02 15:04:05"
	syslogTimeFormat = "2006-01-02T15:04:05.999999Z07:00"
)

var (
	// osExit is a variable for testing purposes.
	osExit = os.Exit

	// coloredLevelPrefixes is used to set the color and format of log prefixes based on log level.
	coloredLevelPrefixes = map[types.Level]string{
		types.DebugLevel:   color.WhiteString(strings.ToUpper(types.LevelToString[types.DebugLevel])),
		types.InfoLevel:    color.GreenString(strings.ToUpper(types.LevelToString[types.InfoLevel])),
		types.WarningLevel: color.YellowString(strings.ToUpper(types.LevelToString[types.WarningLevel])),
		types.ErrorLevel:   color.HiRedString(strings.ToUpper(types.LevelToString[types.ErrorLevel])),
		types.FatalLevel:   color.RedString(strings.ToUpper(types.LevelToString[types.FatalLevel])),
	}
)

// Logger allows messages with different levels/priorities to be sent to stderr/stdout
type Logger struct {
	appName string
	logger  *log.Logger
	format  types.Format
	level   types.Level
	output  func(*types.Entry, int)
}

// New creates a new Logger instance with default configuration settings.
func New() *Logger {
	var config = &types.Config{}
	config.SetDefaults()
	return WithConfig(config)
}

// WithConfig initializes a new Logger with the provided configuration.
// It sets the logger's application name to the basename of the executable.
// If any errors occur during configuration, it panics.
func WithConfig(config *types.Config) *Logger {
	executable, err := os.Executable()
	if err != nil {
		executable = "-"
	}

	l := &Logger{
		appName: path.Base(executable),
		logger:  log.New(os.Stderr, "", 0),
	}

	err = l.SetLevel(config.Level)
	if err != nil {
		panic(err)
	}

	err = l.SetFormat(config.Format)
	if err != nil {
		panic(err)
	}

	err = l.SetOutput(config.Output)
	if err != nil {
		panic(err)
	}

	return l
}

// Format returns the current output format setting of the logger.
func (r *Logger) Format() types.Format {
	return r.format
}

// SetFormat changes the output format of the logger to the specified format.
func (r *Logger) SetFormat(format types.Format) (err error) {
	switch format {
	case types.PlainFormat:
		r.output = r.outputPlain
	case types.JSONFormat:
		r.output = r.outputJSON
	case types.SyslogFormat:
		r.output = r.outputSyslog
	default:
		return types.ErrInvalidFormat
	}
	r.format = format
	return nil
}

// Level returns the current logging level of the logger.
func (r *Logger) Level() types.Level {
	return r.level
}

// SetLevel changes the logging level of the logger to the specified level.
func (r *Logger) SetLevel(level types.Level) error {
	if level < types.DebugLevel || level > types.FatalLevel {
		return types.ErrInvalidLevel
	}
	r.level = level
	return nil
}

// Output is used to log messages with a specific level and format.
func (r *Logger) Output(level types.Level, calldepth int, format *string, v ...any) {
	// If the current logger's level is higher than the given level, return immediately without logging anything.
	if r.level > level {
		return
	}

	var (
		file string
		line int
		ok   bool
		l    *types.Entry
	)

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

	// Create a new Entry with the given level, filename, line number, and message.
	// If format is nil, use fmt.Sprint to create the log message; otherwise, use fmt.Sprintf with the provided format string.
	if format == nil {
		if len(v) == 1 {
			l = types.NewEntry(level, file, line, v[0])
		} else {
			l = types.NewEntry(level, file, line, fmt.Sprint(v...))
		}
	} else {
		l = types.NewEntry(level, file, line, fmt.Sprintf(*format, v...))
	}

	// Pass the Entry to the output function for further processing and logging.
	calldepth++
	r.output(l, calldepth)
}

// outputPlain is a wrapper function to create a log entry in plain text.
func (r *Logger) outputPlain(entry *types.Entry, calldepth int) {
	var b strings.Builder

	b.WriteRune('[')
	b.WriteString(coloredLevelPrefixes[entry.Level])
	b.WriteRune(']')
	for i := 0; i <= 7-len(types.LevelToString[entry.Level]); i++ {
		b.WriteRune(' ')
	}
	b.WriteString(entry.Date.Format(plainTimeFormat))
	b.WriteRune(' ')
	b.WriteString(entry.File)
	b.WriteRune(':')
	b.WriteString(strconv.Itoa(entry.Line))
	b.WriteString(": ")
	b.WriteString(entry.String())

	_ = r.logger.Output(calldepth, b.String())
}

// outputJSON is a wrapper function to create a log entry in JSON.
func (r *Logger) outputJSON(entry *types.Entry, calldepth int) {
	text, err := entry.Marshal()
	if err != nil {
		_ = r.logger.Output(calldepth, err.Error())
	} else {
		_ = r.logger.Output(calldepth, string(text))
	}
}

func (r *Logger) outputSyslog(entry *types.Entry, calldepth int) {
	var b strings.Builder

	b.WriteRune('<')

	// The syslog facility is calculated using the formula: (facility * 8) + severity
	// In this case, we're using facility level "LOCAL0".
	b.WriteString(strconv.FormatInt(int64(16*8+types.LevelToSyslog[entry.Level]), 10))
	b.WriteString(">1 ")
	b.WriteString(entry.Date.Format(syslogTimeFormat))

	// Skip HOSTNAME
	b.WriteString(" - ")
	b.WriteString(r.appName)

	// Skip PROCID and MSGID
	b.WriteString(` - - [goSDID@32473 file="`)
	b.WriteString(entry.File)
	b.WriteString(`" line="`)
	b.WriteString(strconv.Itoa(entry.Line))
	b.WriteString(`"] `)
	b.WriteString(entry.String())

	_ = r.logger.Output(calldepth, b.String())
}

// SetOutput opens or create the given file and set it as the new logging destination.
func (r *Logger) SetOutput(filename string) error {
	file, err := os.OpenFile(path.Clean(filename), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if r.logger != nil {
		r.logger.SetOutput(file)
	}

	return err
}

// Debug writes debug level messages
func (r *Logger) Debug(v ...any) {
	//nolint:gomnd
	r.Output(types.DebugLevel, 2, nil, v...)
}

// Debugf writes debug level messages using formatted string
func (r *Logger) Debugf(format string, v ...any) {
	//nolint:gomnd
	r.Output(types.DebugLevel, 2, &format, v...)
}

// Info writes info level messages
func (r *Logger) Info(v ...any) {
	//nolint:gomnd
	r.Output(types.InfoLevel, 2, nil, v...)
}

// Infof writes info level messages using formatted string
func (r *Logger) Infof(format string, v ...any) {
	//nolint:gomnd
	r.Output(types.InfoLevel, 2, &format, v...)
}

// Warning writes warning level messages
func (r *Logger) Warning(v ...any) {
	//nolint:gomnd
	r.Output(types.WarningLevel, 2, nil, v...)
}

// Warningf writes warning level messages using formatted string
func (r *Logger) Warningf(format string, v ...any) {
	//nolint:gomnd
	r.Output(types.WarningLevel, 2, &format, v...)
}

func (r *Logger) Warn(v ...any) {
	r.Warning(v...)
}

func (r *Logger) Warnf(format string, v ...any) {
	r.Warningf(format, v...)
}

// Error writes error level messages
func (r *Logger) Error(v ...any) {
	//nolint:gomnd
	r.Output(types.ErrorLevel, 2, nil, v...)
}

// Errorf writes error level messages using formatted string
func (r *Logger) Errorf(format string, v ...any) {
	//nolint:gomnd
	r.Output(types.ErrorLevel, 2, &format, v...)
}

// Fatal writes fatal level messages and terminates the application
func (r *Logger) Fatal(v ...any) {
	//nolint:gomnd
	r.Output(types.FatalLevel, 2, nil, v...)
	osExit(1)
}

// Fatalf writes fatal level messages using formatted string and terminates the application
func (r *Logger) Fatalf(format string, v ...any) {
	//nolint:gomnd
	r.Output(types.FatalLevel, 2, &format, v...)
	osExit(1)
}
