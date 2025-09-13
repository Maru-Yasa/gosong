// ...existing code...

package logger

import (
	"os"

	"github.com/charmbracelet/log"
)

type ConsoleLogger struct {
	logger *log.Logger
}

func NewConsoleLogger() *ConsoleLogger {
	return &ConsoleLogger{
		logger: log.NewWithOptions(os.Stdout, log.Options{ReportTimestamp: false}),
	}
}

func (c *ConsoleLogger) Info(format string, args ...interface{}) {
	c.logger.Infof(format, args...)
}

func (c *ConsoleLogger) Warn(format string, args ...interface{}) {
	c.logger.Warnf(format, args...)
}

func (c *ConsoleLogger) Error(format string, args ...interface{}) {
	c.logger.Errorf(format, args...)
}
