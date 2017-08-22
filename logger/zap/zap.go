package zap

import (
	"github.com/cryptopay-dev/gemstone/logger"
	"go.uber.org/zap"
)

type Logger struct {
	logger *zap.SugaredLogger
}

func New() logger.Logger {
	logger, _ := zap.NewProduction(zap.AddCallerSkip(1))
	sugar := logger.Sugar()

	return &Logger{
		logger: sugar,
	}
}

func (logger *Logger) Named(name string) logger.Logger {
	return &Logger{
		logger: logger.logger.Named(name),
	}
}

func (logger *Logger) WithContext(args ...interface{}) logger.Logger {
	return &Logger{
		logger: logger.logger.With(args...),
	}
}

func (logger *Logger) Info(args ...interface{}) {
	logger.logger.Info(args...)
}

func (logger *Logger) Debug(args ...interface{}) {
	logger.logger.Info(args...)
}

func (logger *Logger) Warning(args ...interface{}) {
	logger.logger.Warn(args...)
}

func (logger *Logger) Error(args ...interface{}) {
	logger.logger.Error(args...)
}

func (logger *Logger) Panic(args ...interface{}) {
	logger.logger.Panic(args...)
}

// With external params
func (logger *Logger) Infow(message string, args ...interface{}) {
	logger.logger.Infow(message, args...)
}

func (logger *Logger) Debugw(message string, args ...interface{}) {
	logger.logger.Debugw(message, args...)
}

func (logger *Logger) Warningw(message string, args ...interface{}) {
	logger.logger.Warnw(message, args...)
}

func (logger *Logger) Errorw(message string, args ...interface{}) {
	logger.logger.Errorw(message, args...)
}

func (logger *Logger) Panicw(message string, args ...interface{}) {
	logger.logger.Panicw(message, args...)
}

// Using formatting
func (logger *Logger) Infof(message string, args ...interface{}) {
	logger.logger.Infof(message, args...)
}

func (logger *Logger) Debugf(message string, args ...interface{}) {
	logger.logger.Debugf(message, args...)
}

func (logger *Logger) Warningf(message string, args ...interface{}) {
	logger.logger.Warnf(message, args...)
}

func (logger *Logger) Errorf(message string, args ...interface{}) {
	logger.logger.Errorf(message, args...)
}

func (logger *Logger) Panicf(message string, args ...interface{}) {
	logger.logger.Panicf(message, args...)
}
