package logging

type Logger interface {
	Info(message string)
	Infof(format string, params ...any)
	Warn(message string)
	Warnf(format string, params ...any)
	Debug(message string)
	Debugf(format string, params ...any)
	Error(err error, message string)
	Errorf(err error, format string, params ...any)
	Fatal(err error, message string)
	Fatalf(err error, format string, params ...any)
}
