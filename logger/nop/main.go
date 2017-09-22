package nop

import (
	"github.com/cryptopay-dev/gemstone/logger"
)

type Logger struct{}

func New() logger.Logger                                        { return &Logger{} }
func (l *Logger) WithContext(args ...interface{}) logger.Logger { return l }
func (l *Logger) Named(name string) logger.Logger               { return l }
func (l *Logger) Info(args ...interface{})                      {}
func (l *Logger) Debug(args ...interface{})                     {}
func (l *Logger) Warning(args ...interface{})                   {}
func (l *Logger) Error(args ...interface{})                     {}
func (l *Logger) Panic(args ...interface{})                     {}
func (l *Logger) Infow(message string, args ...interface{})     {}
func (l *Logger) Debugw(message string, args ...interface{})    {}
func (l *Logger) Warningw(message string, args ...interface{})  {}
func (l *Logger) Errorw(message string, args ...interface{})    {}
func (l *Logger) Panicw(message string, args ...interface{})    {}
func (l *Logger) Infof(message string, args ...interface{})     {}
func (l *Logger) Debugf(message string, args ...interface{})    {}
func (l *Logger) Warningf(message string, args ...interface{})  {}
func (l *Logger) Errorf(message string, args ...interface{})    {}
func (l *Logger) Panicf(message string, args ...interface{})    {}
