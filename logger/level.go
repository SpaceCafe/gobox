package logger

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

	// WarningLevel messages are used when something unexpected happened but did not cause the program to fail.
	// For example, "Failed to connect to database", "Invalid configuration value".
	// It's important to monitor these logs as they often indicate a problem that may need attention soon.
	WarningLevel

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

	// WarnLevel is an alias of WarningLevel.
	WarnLevel = WarningLevel
)

var (
	// LevelToString is a map that converts a Level to its string representation.
	LevelToString = map[Level]string{DebugLevel: "debug", InfoLevel: "info", WarningLevel: "warning", ErrorLevel: "error", FatalLevel: "fatal"}

	// StringToLevel is a map that converts a string to its Level equivalent.
	StringToLevel = map[string]Level{"debug": DebugLevel, "info": InfoLevel, "warning": WarningLevel, "error": ErrorLevel, "fatal": FatalLevel}
)
