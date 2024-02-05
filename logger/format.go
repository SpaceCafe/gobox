package logger

// Format is used to specify the format in which logs should be written.
type Format int

const (
	// PlainFormat represents plain text logs (default).
	PlainFormat Format = 0 + iota

	// JSONFormat represents logs in JSON format.
	JSONFormat
)
