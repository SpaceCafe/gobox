package logger

import (
	"encoding/json"
	"fmt"
	"time"
)

// Entry represents a log item with date, file name, line number, level, and Message.
type Entry struct {
	// Date is the timestamp of the log.
	Date time.Time `json:"date"`

	// File represents the source code file where the log was created.
	File string `json:"file,omitempty"`

	// Level is the logging level (e.g., "info", "warning").
	Level Level `json:"level"`

	// Message contains the log entry.
	// It could be a string or a fmt.Stringer interface with JSON annotations.
	Message any `json:"message"`

	// Line contains the corresponding line number where the log was created.
	Line int `json:"line,omitempty"`
}

// NewEntry creates a new Entry with the given parameters and current time.
func NewEntry(level Level, file string, line int, message any) *Entry {
	return &Entry{
		Date:    time.Now(),
		File:    file,
		Level:   level,
		Message: message,
		Line:    line,
	}
}

// Marshal converts an Entry to a JSON byte array and returns it.
// If there's an error during conversion, it wraps the error.
func (r *Entry) Marshal() ([]byte, error) {
	if r.Message == nil {
		r.Message = ""
	}

	out, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Entry: %w", err)
	}

	return out, nil
}

// String returns the string representation of the Entry's Message.
// It handles nil, string, fmt.Stringer, and other types accordingly.
func (r *Entry) String() string {
	switch msg := r.Message.(type) {
	case nil:
		return ""
	case string:
		return msg
	case fmt.Stringer:
		return msg.String()
	default:
		return fmt.Sprintf("%v", msg)
	}
}
