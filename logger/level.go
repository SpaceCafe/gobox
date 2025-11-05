package logger

import (
	"errors"
)

type Level int

const (
	// DebugLevel is the lowest level and usually only enabled during development.
	// It provides detailed information useful for debugging purposes,
	// such as tracking variable values or understanding the flow of execution through a program.
	// In production environments, this level should be disabled to reduce log size and improve performance.
	DebugLevel Level = 0 + iota

	// InfoLevel messages are typically used when you want to provide contextual details
	// about normal operations in your application.
	// For example, "Starting server on port 8080" or "Connected to database".
	// These logs can help with monitoring the health of an application and understanding its behavior over time.
	InfoLevel

	// WarnLevel messages are used when something unexpected happened but did not cause the program to fail.
	// For example, "Failed to connect to database", "Invalid configuration value".
	// It's important to monitor these logs as they often indicate a problem that may need attention soon.
	WarnLevel

	// ErrorLevel messages are for reporting failures in a way that allows the application to continue running.
	// For instance, "Could not find file", "Failed to parse JSON".
	// These types of errors should be fixed immediately but do not necessarily mean the program is failing.
	ErrorLevel

	// FatalLevel logs indicate that an unrecoverable error has occurred and the application cannot proceed
	// with its current operation. The application will typically exit after logging a fatal message,
	// possibly with a non-zero status code to indicate failure.
	// For example, "Failed to start server", "Database connection lost".
	// These types of errors are usually indicative of serious problems that require immediate attention.
	FatalLevel
)

var (
	ErrInvalidLevel = errors.New("log level is invalid")

	// LevelToString is a map that converts a Level to its string representation.
	//nolint:gochecknoglobals // This is a lookup map that needs to be globally accessible.
	LevelToString = map[Level]string{
		DebugLevel: "debug",
		InfoLevel:  "info",
		WarnLevel:  "warn",
		ErrorLevel: "error",
		FatalLevel: "fatal",
	}

	// LevelToSyslog maps log levels to their corresponding syslog severity values.
	//nolint:gochecknoglobals // This is a lookup map that needs to be globally accessible.
	LevelToSyslog = map[Level]int{
		//nolint:mnd // These numbers are standard syslog severity levels defined in RFC 5424.
		DebugLevel: 7, InfoLevel: 6, WarnLevel: 4, ErrorLevel: 3, FatalLevel: 2,
	}

	// StringToLevel is a map that converts a string to its Level equivalent.
	//nolint:gochecknoglobals // This is a lookup map that needs to be globally accessible.
	StringToLevel = map[string]Level{
		"debug":   DebugLevel,
		"info":    InfoLevel,
		"warn":    WarnLevel,
		"warning": WarnLevel,
		"error":   ErrorLevel,
		"fatal":   FatalLevel,
	}
)

// MarshalText serializes the Level to a textual representation.
func (r *Level) MarshalText() ([]byte, error) {
	return []byte(r.String()), nil
}

// String returns the string representation of the Level.
func (r *Level) String() string {
	return LevelToString[*r]
}

// UnmarshalText converts a textual representation of the log level into a Level type.
func (r *Level) UnmarshalText(text []byte) (err error) {
	*r, err = ParseLevel(string(text))

	return
}

// ParseLevel converts a string to its corresponding Level type.
// If the conversion fails, it returns the InfoLevel and an error.
func ParseLevel(level string) (Level, error) {
	if v, ok := StringToLevel[level]; ok {
		return v, nil
	}

	return InfoLevel, ErrInvalidLevel
}
