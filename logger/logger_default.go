package logger

// std is the default logger.
var std = New()

// Default returns the default logger.
func Default() *Logger {
	return std
}

// Aliases of logger functions.
var (
	SetFormat  = std.SetFormat
	SetLevel   = std.SetLevel
	ParseLevel = std.ParseLevel
	SetOutput  = std.SetOutput
	Output     = std.Output
	Debug      = std.Debug
	Debugf     = std.Debugf
	Info       = std.Info
	Infof      = std.Infof
	Warning    = std.Warning
	Warningf   = std.Warningf
	Warn       = std.Warn
	Warnf      = std.Warnf
	Error      = std.Error
	Errorf     = std.Errorf
	Fatal      = std.Fatal
	Fatalf     = std.Fatalf
)
