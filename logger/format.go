package logger

import (
	"errors"
)

// Format is used to specify the format in which logs should be written.
type Format int

const (
	// PlainFormat represents plain text logs (default).
	PlainFormat Format = 0 + iota

	// JSONFormat represents logs in JSON format.
	JSONFormat

	// SyslogFormat represents logs in syslog format.
	SyslogFormat
)

var (
	ErrInvalidFormat = errors.New("log format is invalid")

	// FormatToString is a map that converts a Format to its string representation.
	//nolint:gochecknoglobals // This is a lookup map that needs to be globally accessible.
	FormatToString = map[Format]string{
		PlainFormat:  "plain",
		JSONFormat:   "json",
		SyslogFormat: "syslog",
	}

	// StringToFormat is a map that converts a string to its Format equivalent.
	//nolint:gochecknoglobals // This is a lookup map that needs to be globally accessible.
	StringToFormat = map[string]Format{
		"plain":  PlainFormat,
		"json":   JSONFormat,
		"syslog": SyslogFormat,
	}
)

func (r *Format) String() string {
	return FormatToString[*r]
}

func (r *Format) MarshalText() ([]byte, error) {
	return []byte(r.String()), nil
}

func (r *Format) UnmarshalText(text []byte) (err error) {
	*r, err = ParseFormat(string(text))
	return
}

func ParseFormat(format string) (Format, error) {
	if v, ok := StringToFormat[format]; ok {
		return v, nil
	}
	return PlainFormat, ErrInvalidFormat
}
