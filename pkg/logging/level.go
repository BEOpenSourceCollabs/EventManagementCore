package logging

type Level int8

const (
	INFO Level = iota
	WARN
	DEBUG
	ERROR
	FATAL
	OFF
)

func (l Level) String() string {
	switch l {
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case DEBUG:
		return "DEBUG"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return ""
	}
}
