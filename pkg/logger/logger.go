package logger

import (
	"sync"
)

// Logger is the interface for logging
type Logger interface {
	Info(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Error(format string, args ...interface{})
}

// MultiLogger broadcasts logs to multiple loggers
type MultiLogger struct {
	loggers []Logger
}

func NewMultiLogger(loggers ...Logger) *MultiLogger {
	return &MultiLogger{loggers: loggers}
}

func (m *MultiLogger) Info(format string, args ...interface{}) {
	for _, l := range m.loggers {
		l.Info(format, args...)
	}
}

func (m *MultiLogger) Warn(format string, args ...interface{}) {
	for _, l := range m.loggers {
		l.Warn(format, args...)
	}
}

func (m *MultiLogger) Error(format string, args ...interface{}) {
	for _, l := range m.loggers {
		l.Error(format, args...)
	}
}

var (
	defaultLogger Logger = &noopLogger{}
	once          sync.Once
)

// SetDefaultLogger sets the global logger (should be called once, e.g. in main)
func SetDefaultLogger(l Logger) {
	once.Do(func() {
		defaultLogger = l
	})
}

// Info logs info message to the default logger
func Info(format string, args ...interface{}) {
	defaultLogger.Info(format, args...)
}

// Warn logs warning message to the default logger
func Warn(format string, args ...interface{}) {
	defaultLogger.Warn(format, args...)
}

// Error logs error message to the default logger
func Error(format string, args ...interface{}) {
	defaultLogger.Error(format, args...)
}

// noopLogger implements Logger but does nothing (for zero value safety)
type noopLogger struct{}

func (n *noopLogger) Info(string, ...interface{})  {}
func (n *noopLogger) Warn(string, ...interface{})  {}
func (n *noopLogger) Error(string, ...interface{}) {}
