package logger

// Logger
type Logger interface {
	Infof(format string, args ...interface{})

	Warnf(format string, args ...interface{})

	Errorf(format string, args ...interface{})

	WithFields(keyValues Fields) Logger

	WithError(err error) Logger
}

// Fields
type Fields map[string]interface{}

// LoggerAdapter
type LoggerAdapter struct {
	log Logger
}

// NewLoggerAdapter create new LoggerAdapter with its dependencies
func NewLoggerAdapter(adapter Logger) *LoggerAdapter {
	return &LoggerAdapter{log: adapter}
}

// Log return the log property
func (a LoggerAdapter) Log() Logger {
	return a.log
}
