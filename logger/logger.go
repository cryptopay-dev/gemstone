package logger

type Logger interface {
	WithContext(args ...interface{}) Logger

	Info(args ...interface{})
	Debug(args ...interface{})
	Warning(args ...interface{})
	Erorr(args ...interface{})
	Panic(args ...interface{})

	Infow(message string, args ...interface{})
	Debugw(message string, args ...interface{})
	Warningw(message string, args ...interface{})
	Erorrw(message string, args ...interface{})
	Panicw(message string, args ...interface{})

	Infof(message string, args ...interface{})
	Debugf(message string, args ...interface{})
	Warningf(message string, args ...interface{})
	Erorrf(message string, args ...interface{})
	Panicf(message string, args ...interface{})
}
