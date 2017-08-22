package logger

import "go.uber.org/zap"

type ZapLogger struct {
	logger *zap.SugaredLogger
}

func NewZap() Logger {
	logger, _ := zap.NewProduction()
	sugar := logger.Sugar()

	return &ZapLogger{
		logger: sugar,
	}
}

func (logger *ZapLogger) WithContext(args ...interface{}) Logger {
	zap := logger.logger.With(args...)

	return &ZapLogger{
		logger: zap,
	}
}

func (logger *ZapLogger) Info(args ...interface{}) {
	logger.logger.Info(args...)
}

func (logger *ZapLogger) Debug(args ...interface{}) {
	logger.logger.Info(args...)
}

func (logger *ZapLogger) Warning(args ...interface{}) {
	logger.logger.Warn(args...)
}

func (logger *ZapLogger) Erorr(args ...interface{}) {
	logger.logger.Error(args...)
}

func (logger *ZapLogger) Panic(args ...interface{}) {
	logger.logger.Panic(args...)
}

// With external params
func (logger *ZapLogger) Infow(message string, args ...interface{}) {
	logger.logger.Infow(message, args...)
}

func (logger *ZapLogger) Debugw(message string, args ...interface{}) {
	logger.logger.Debugw(message, args...)
}

func (logger *ZapLogger) Warningw(message string, args ...interface{}) {
	logger.logger.Warnw(message, args...)
}

func (logger *ZapLogger) Erorrw(message string, args ...interface{}) {
	logger.logger.Errorw(message, args...)
}

func (logger *ZapLogger) Panicw(message string, args ...interface{}) {
	logger.logger.Panicw(message, args...)
}

// Using formatting
func (logger *ZapLogger) Infof(message string, args ...interface{}) {
	logger.logger.Infof(message, args...)
}

func (logger *ZapLogger) Debugf(message string, args ...interface{}) {
	logger.logger.Debugf(message, args...)
}

func (logger *ZapLogger) Warningf(message string, args ...interface{}) {
	logger.logger.Warnf(message, args...)
}

func (logger *ZapLogger) Erorrf(message string, args ...interface{}) {
	logger.logger.Errorf(message, args...)
}

func (logger *ZapLogger) Panicf(message string, args ...interface{}) {
	logger.logger.Panicf(message, args...)
}
