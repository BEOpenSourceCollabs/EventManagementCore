package logging

import (
	"fmt"
	"os"
)

type ContextLogger struct {
	lw      LogWriter
	Context string
}

func NewContextLogger(lw LogWriter, context string) *ContextLogger {
	return &ContextLogger{
		lw:      lw,
		Context: context,
	}
}

func (cl *ContextLogger) Info(message string) {
	cl.lw.print(INFO, message, cl.Context)
}

func (cl *ContextLogger) Infof(format string, params ...any) {
	cl.lw.print(INFO, fmt.Sprintf(format, params...), cl.Context)
}

func (cl *ContextLogger) Warn(message string) {
	cl.lw.print(WARN, message, cl.Context)
}

func (cl *ContextLogger) Warnf(format string, params ...any) {
	cl.lw.print(WARN, fmt.Sprintf(format, params...), cl.Context)
}

func (cl *ContextLogger) Debug(message string) {
	cl.lw.print(DEBUG, message, cl.Context)
}

func (cl *ContextLogger) Debugf(format string, params ...any) {
	cl.lw.print(DEBUG, fmt.Sprintf(format, params...), cl.Context)
}

func (cl *ContextLogger) Error(err error, message string) {
	cl.lw.print(ERROR, fmt.Sprintf("(%s) %s", err.Error(), message), cl.Context)
}

func (cl *ContextLogger) Errorf(err error, format string, params ...any) {
	cl.lw.print(ERROR, fmt.Sprintf("(%s) %s", err.Error(), fmt.Sprintf(format, params...)), cl.Context)
}

func (cl *ContextLogger) Fatal(err error, message string) {
	cl.lw.print(FATAL, fmt.Sprintf("(%s) %s", err.Error(), message), cl.Context)
	os.Exit(1)
}

func (cl *ContextLogger) Fatalf(err error, format string, params ...any) {
	cl.lw.print(FATAL, fmt.Sprintf("%s: %s", fmt.Sprintf(format, params...), err.Error()), cl.Context)
	os.Exit(1)
}
