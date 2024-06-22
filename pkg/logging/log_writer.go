package logging

type LogWriter interface {
	print(level Level, message string, context string) (int, error)
}
