package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type FileLogger struct {
	logger *zap.SugaredLogger
}

func NewFileLogger() *FileLogger {
	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewJSONEncoder(cfg)
	file, _ := os.OpenFile("gosong.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	core := zapcore.NewCore(encoder, zapcore.AddSync(file), zapcore.InfoLevel)
	logger := zap.New(core).Sugar()
	return &FileLogger{logger: logger}
}

func (f *FileLogger) Info(format string, args ...interface{}) {
	f.logger.Infof(format, args...)
}

func (f *FileLogger) Warn(format string, args ...interface{}) {
	f.logger.Warnf(format, args...)
}

func (f *FileLogger) Error(format string, args ...interface{}) {
	f.logger.Errorf(format, args...)
}

func (f *FileLogger) Sync() error {
	return f.logger.Sync()
}
