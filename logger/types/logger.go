package types

type Logger interface {
	Level() Level
	SetLevel(Level) error

	Format() Format
	SetFormat(Format) error

	Output(level Level, calldepth int, format *string, v ...any)
	SetOutput(filename string) error

	Debug(v ...any)
	Debugf(format string, v ...any)
	Info(v ...any)
	Infof(format string, v ...any)
	Warn(v ...any)
	Warnf(format string, v ...any)
	Error(v ...any)
	Errorf(format string, v ...any)
	Fatal(v ...any)
	Fatalf(format string, v ...any)
}
