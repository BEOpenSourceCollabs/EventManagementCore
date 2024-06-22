package logging

type Logger interface {
	Info(context, message string)
	Infof(context, format string, params ...any)
	Warn(context, message string)
	Warnf(context, format string, params ...any)
	Debug(context, message string)
	Debugf(context, format string, params ...any)
	Error(context string, err error, message string)
	Errorf(context string, err error, format string, params ...any)
	Fatal(context string, err error, message string)
	Fatalf(context string, err error, format string, params ...any)
}
