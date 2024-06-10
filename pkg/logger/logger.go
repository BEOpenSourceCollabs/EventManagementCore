package logger

import (
	"fmt"
	"log"
	"os"
)

var (
	AppLogger *appLogger
)

type appLogger struct {
	defaultLogger *log.Logger
}

func init() {
	AppLogger = &appLogger{
		defaultLogger: log.New(os.Stdout, "", log.Ldate|log.Ltime),
	}
}

// prints log of type info with message and context
func (al *appLogger) Info(context, message string) {
	al.defaultLogger.Printf("INFO[%s]: %s", context, message)
}

// prints log of type info and pass formatted string and context
func (al *appLogger) InfoF(context, format string, v ...any) {
	al.defaultLogger.Printf("INFO[%s]: %s", context, fmt.Sprintf(format, v...))
}

// prints log of type WARN with message and contex
func (al *appLogger) Warn(context, message string) {
	al.defaultLogger.Printf("WARN[%s]: %s", context, message)
}

// prints log of type WARN and pass formatted string and context
func (al *appLogger) WarnF(context, format string, v ...any) {
	al.defaultLogger.Printf("WARN[%s]: %s", context, fmt.Sprintf(format, v...))
}

// prints log of type ERROR with message and context
func (al *appLogger) Error(context, message string) {
	al.defaultLogger.Printf("ERROR[%s]: %s", context, message)
}

// prints log of type ERROR and pass formatted string and context
func (al *appLogger) ErrorF(context, format string, v ...any) {
	al.defaultLogger.Printf("ERROR[%s]: %s", context, fmt.Sprintf(format, v...))
}

func (al *appLogger) Fatal(context string, err error) {
	al.defaultLogger.Fatalf("Fatal[%s]: %v", context, err)
}
