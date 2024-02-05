package logger

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

const (
	// prefixMaxLen sets the maximum length of log prefixes.
	prefixMaxLen = 7

	// dateFormat sets the format using Go's reference time.
	dateFormat = "2006/01/02 15:04:05"
)

var (
	// prefix is used to set the color and format of log prefixes based on log level.
	prefix = map[Level]string{
		DebugLevel:   color.WhiteString(strings.ToUpper(LevelToString[DebugLevel])),
		InfoLevel:    color.GreenString(strings.ToUpper(LevelToString[InfoLevel])),
		WarningLevel: color.YellowString(strings.ToUpper(LevelToString[WarningLevel])),
		ErrorLevel:   color.HiRedString(strings.ToUpper(LevelToString[ErrorLevel])),
		FatalLevel:   color.RedString(strings.ToUpper(LevelToString[FatalLevel])),
	}
)

// Item represents a log item with date, file name, line number, level, and Message.
type Item struct {
	// Date is the timestamp of the log.
	Date time.Time `json:"date"`

	// File represents the source code file where the log was created.
	File string `json:"file"`

	// Level is a string representation of the logging level (e.g., "info", "warning").
	Level string `json:"level"`

	// Message contains the log entry.
	// It could be a string or a fmt.Stringer interface with json annotations.
	Message any `json:"message"`

	// Line contains the corresponding line number where the log was created.
	Line int `json:"line"`

	// level is actual logging level of the log.
	level Level
}

// NewItem creates a new Item with the given parameters and current time.
func NewItem(level Level, file string, line int, message any) *Item {
	return &Item{
		Date:    time.Now(),
		File:    file,
		Level:   LevelToString[level],
		Message: message,
		Line:    line,
		level:   level,
	}
}

// Marshal converts an Item to a JSON byte array and returns it.
// If there's an error during conversion, it wraps the error.
func (r *Item) Marshal() ([]byte, error) {
	if r.Message == nil {
		r.Message = ""
	}

	out, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Item: %w", err)
	}
	return out, nil
}

// String returns a string representation of the log item in the format "[LEVEL] [DATE] [FILE]:[LINE]: MESSAGE".
// The level is colored based on its value.
func (r *Item) String() string {
	buf := "[" + prefix[r.level] + "]"
	for i := 0; i <= prefixMaxLen-len(prefix[r.level]); i++ {
		buf += " "
	}
	buf += r.Date.Format(dateFormat) + " " + r.File + ":" + strconv.Itoa(r.Line) + ": "

	switch msg := r.Message.(type) {
	case nil:
		break
	case string:
		buf += msg
	default:
		buf += fmt.Sprintf("%v", msg)
	}

	return buf
}
