// Package logging provides utilities for logging.
package logging

type Config struct {
	Level Level `default:"info"`
}

// Fields type, used to pass to `WithFields`.
type Fields map[string]interface{}

type Level string

var (
	ErrorLevel = Level("error")
	InfoLevel  = Level("info")
	WarnLevel  = Level("warn")
	DebugLevel = Level("debug")
	TraceLevel = Level("trace")
)

// Logger defines the generic logging interface
// https://i.stack.imgur.com/z5Fim.png
type Logger interface {
	Infof(format string, args ...interface{})
	Info(args ...interface{})

	Errorf(format string, args ...interface{})
	Error(args ...interface{})

	Warnf(format string, args ...interface{})
	Warn(args ...interface{})

	Debugf(format string, args ...interface{})
	Debug(args ...interface{})

	Tracef(format string, args ...interface{})
	Trace(args ...interface{})

	WithFields(fields Fields) Logger
	WithField(field string, arg interface{}) Logger

	WithError(err error) Logger
}
