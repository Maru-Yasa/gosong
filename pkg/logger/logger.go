package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.SugaredLogger

func Init() error {
	cfg := zap.NewDevelopmentEncoderConfig()
	cfg.EncodeTime = zapcore.TimeEncoderOfLayout("2006/01/02 15:04:05")
	encoder := zapcore.NewConsoleEncoder(cfg)
	core := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.InfoLevel)
	logger := zap.New(core)
	log = logger.Sugar()
	return nil
}

func Info(msg string, host string, fields ...interface{}) {
	prefix := fmt.Sprintf("[%s]", host)
	log.Infof("%s %s", prefix, fmt.Sprintf(msg, fields...))
}

func Error(msg string, host string, fields ...interface{}) {
	prefix := fmt.Sprintf("[%s]", host)
	log.Errorf("%s %s", prefix, fmt.Sprintf(msg, fields...))
}

func Sync() error {
	if log != nil {
		return log.Sync()
	}
	return nil
}
