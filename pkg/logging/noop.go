package logging

// NoopLogger does not do anything. Useful for testing.
type NoopLogger struct{}

func (NoopLogger) Infof(format string, args ...interface{})  {}
func (NoopLogger) Info(args ...interface{})                  {}
func (NoopLogger) Errorf(format string, args ...interface{}) {}
func (NoopLogger) Error(args ...interface{})                 {}
func (NoopLogger) Warnf(format string, args ...interface{})  {}
func (NoopLogger) Warn(args ...interface{})                  {}
func (NoopLogger) Debugf(format string, args ...interface{}) {}
func (NoopLogger) Debug(args ...interface{})                 {}
func (NoopLogger) Tracef(format string, args ...interface{}) {}
func (NoopLogger) Trace(args ...interface{})                 {}

func (n NoopLogger) WithFields(fields Fields) Logger                { return n }
func (n NoopLogger) WithField(field string, arg interface{}) Logger { return n }
func (n NoopLogger) WithError(err error) Logger                     { return n }
