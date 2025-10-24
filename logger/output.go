package logger

import (
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

const (
	// plainTimeFormat sets the format using Go's reference time.
	plainTimeFormat  = "2006-01-02 15:04:05"
	syslogTimeFormat = "2006-01-02T15:04:05.999999Z07:00"
)

var (
	// osExit is a variable for testing purposes.
	//nolint:gochecknoglobals // This is a mock for os.Exit used in tests to prevent actual program termination
	osExit = os.Exit

	// coloredLevelPrefixes is used to set the color and format of log prefixes based on log level.
	//nolint:gochecknoglobals // This is a lookup map that needs to be globally accessible.
	coloredLevelPrefixes = map[Level]string{
		DebugLevel: color.WhiteString(strings.ToUpper(LevelToString[DebugLevel])),
		InfoLevel:  color.GreenString(strings.ToUpper(LevelToString[InfoLevel])),
		WarnLevel:  color.YellowString(strings.ToUpper(LevelToString[WarnLevel])),
		ErrorLevel: color.HiRedString(strings.ToUpper(LevelToString[ErrorLevel])),
		FatalLevel: color.RedString(strings.ToUpper(LevelToString[FatalLevel])),
	}
)

// outputPlain is a wrapper function to create a log entry in plain text.
func (r *DefaultLogger) outputPlain(entry *Entry, calldepth int) {
	var b strings.Builder

	b.WriteByte('[')
	b.WriteString(coloredLevelPrefixes[entry.Level])
	b.WriteByte(']')
	for i := 0; i <= 7-len(LevelToString[entry.Level]); i++ {
		b.WriteByte(' ')
	}
	b.WriteString(entry.Date.Format(plainTimeFormat))
	b.WriteByte(' ')
	b.WriteString(entry.File)
	b.WriteByte(':')
	b.WriteString(strconv.Itoa(entry.Line))
	b.WriteString(": ")
	b.WriteString(entry.String())

	_ = r.logger.Output(calldepth, b.String())
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
	var b strings.Builder

	b.WriteByte('<')

	// The syslog facility is calculated using the formula: (facility * 8) + severity
	// In this case, we're using facility level "LOCAL0".
	b.WriteString(strconv.FormatInt(int64(16*8+LevelToSyslog[entry.Level]), 10))
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
	b.WriteString(escapeSyslogMessage(entry.String()))

	_ = r.logger.Output(calldepth, b.String())
}

// escapeSyslogMessage escapes special characters in syslog messages
// to prevent log injection and ensure RFC 5424 compliance.
func escapeSyslogMessage(msg string) string {
	replacer := strings.NewReplacer(
		"\n", "\\n",
		"\r", "\\r",
		"\t", "\\t",
		"\\", "\\\\",
	)
	return replacer.Replace(msg)
}
