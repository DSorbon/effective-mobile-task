package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var globalLogger *zap.SugaredLogger

func Init(core zapcore.Core, options ...zap.Option) {
	globalLogger = zap.New(core, options...).Sugar()
}

func Debug(args ...interface{}) {
	globalLogger.Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	globalLogger.Debugf(template, args...)
}

func Info(args ...interface{}) {
	globalLogger.Info(args...)
}

func Infof(template string, args ...interface{}) {
	globalLogger.Infof(template, args...)
}

func Warn(args ...interface{}) {
	globalLogger.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	globalLogger.Warnf(template, args...)
}

func Error(args ...interface{}) {
	globalLogger.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	globalLogger.Errorf(template, args...)
}
