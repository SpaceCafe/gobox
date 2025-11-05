package logger

import (
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

const (
	// PlainTimeFormat sets the format using Go's reference time.
	PlainTimeFormat  = "2006-01-02 15:04:05"
	SyslogTimeFormat = "2006-01-02T15:04:05.999999Z07:00"
)

var (
	// OsExit is a variable for testing purposes.
	//nolint:gochecknoglobals // This is a mock for os.Exit used in tests to prevent actual program termination
	OsExit = os.Exit

	// ColoredLevelPrefixes is used to set the color and format of log prefixes based on log level.
	//nolint:gochecknoglobals // This is a lookup map that needs to be globally accessible.
	ColoredLevelPrefixes = map[Level]string{
		DebugLevel: color.WhiteString(strings.ToUpper(LevelToString[DebugLevel])),
		InfoLevel:  color.GreenString(strings.ToUpper(LevelToString[InfoLevel])),
		WarnLevel:  color.YellowString(strings.ToUpper(LevelToString[WarnLevel])),
		ErrorLevel: color.HiRedString(strings.ToUpper(LevelToString[ErrorLevel])),
		FatalLevel: color.RedString(strings.ToUpper(LevelToString[FatalLevel])),
	}
)

// outputPlain is a wrapper function to create a log entry in plain text.
func (r *DefaultLogger) outputPlain(entry *Entry, calldepth int) {
	var builder strings.Builder

	builder.WriteByte('[')
	builder.WriteString(ColoredLevelPrefixes[entry.Level])
	builder.WriteByte(']')

	for i := len(LevelToString[entry.Level]); i < 7; i++ {
		builder.WriteByte(' ')
	}

	builder.WriteString(entry.Date.Format(PlainTimeFormat))

	if entry.File != "" {
		builder.WriteByte(' ')
		builder.WriteString(entry.File)
		builder.WriteByte(':')
		builder.WriteString(strconv.Itoa(entry.Line))
	}

	builder.WriteString(": ")
	builder.WriteString(entry.String())

	_ = r.logger.Output(calldepth, builder.String())
}

// outputJSON is a wrapper function to create a log entry in JSON.
func (r *DefaultLogger) outputJSON(entry *Entry, calldepth int) {
	text, err := entry.Marshal()
	if err != nil {
		_ = r.logger.Output(calldepth, err.Error())
	} else {
		_ = r.logger.Output(calldepth, string(text))
	}
}

func (r *DefaultLogger) outputSyslog(entry *Entry, calldepth int) {
	var builder strings.Builder

	builder.WriteByte('<')

	// The syslog facility is calculated using the formula: (facility * 8) + severity
	// In this case, we're using facility level "LOCAL0".
	builder.WriteString(strconv.FormatInt(int64(16*8+LevelToSyslog[entry.Level]), 10))
	builder.WriteString(">1 ")
	builder.WriteString(entry.Date.Format(SyslogTimeFormat))

	// Skip HOSTNAME
	builder.WriteString(" - ")
	builder.WriteString(r.appName)

	// Skip PROCID and MSGID
	builder.WriteString(" - - ")

	if entry.File != "" {
		builder.WriteString(`[goSDID@32473 file="`)
		builder.WriteString(entry.File)
		builder.WriteString(`" line="`)
		builder.WriteString(strconv.Itoa(entry.Line))
		builder.WriteString(`"] `)
	} else {
		builder.WriteString("- ")
	}

	builder.WriteString(EscapeSyslogMessage(entry.String()))

	_ = r.logger.Output(calldepth, builder.String())
}

// EscapeSyslogMessage escapes special characters in syslog messages
// to prevent log injection and ensure RFC 5424 compliance.
func EscapeSyslogMessage(msg string) string {
	replacer := strings.NewReplacer(
		"\n", "\\n",
		"\r", "\\r",
		"\t", "\\t",
		"\\", "\\\\",
	)

	return replacer.Replace(msg)
}

// LeftPadString appends the text to the builder, left-padded with spaces to reach the specified length.
func LeftPadString(builder *strings.Builder, text string, length int) {
	if len(text) > length {
		builder.WriteString(text[:length])
	} else {
		builder.WriteString(text)

		for i := len(text); i < length; i++ {
			builder.WriteByte(' ')
		}
	}
}

// RightPadString appends the text to the builder, right-padded with spaces to reach the specified length.
func RightPadString(builder *strings.Builder, text string, length int) {
	if len(text) > length {
		builder.WriteString(text[:length])
	} else {
		for i := len(text); i < length; i++ {
			builder.WriteByte(' ')
		}

		builder.WriteString(text)
	}
}
