package logging

import (
	"fmt"
	"io"
	"runtime/debug"
	"sync"
	"time"
)

type TextLogWriter struct {
	o    io.Writer
	minL Level
	mu   sync.Mutex
}

func NewTextLogWriter(out io.Writer, minLevel Level) *TextLogWriter {
	return &TextLogWriter{
		o:    out,
		minL: minLevel,
	}
}

func (tl *TextLogWriter) print(level Level, message string, context string) (int, error) {

	if level < tl.minL {
		return 0, nil
	}

	t := time.Now().Format(time.RFC3339)
	line := fmt.Sprintf("%s %v[%s]: %s\n", t, level, context, message)

	if level >= ERROR {
		line += fmt.Sprintf("trace: %s", string(debug.Stack()))
	}

	tl.mu.Lock()

	defer tl.mu.Unlock()

	return tl.o.Write([]byte(line))
}
