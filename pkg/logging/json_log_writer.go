package logging

import (
	"encoding/json"
	"fmt"
	"io"
	"runtime/debug"
	"sync"
	"time"
)

type JsonLogWriter struct {
	o    io.Writer
	minL Level
	mu   sync.Mutex
}

func NewJsonLogWriter(out io.Writer, minLevel Level) *JsonLogWriter {
	return &JsonLogWriter{
		minL: minLevel,
		o:    out,
	}
}

func (jl *JsonLogWriter) print(level Level, message string, context string) (int, error) {

	if level < jl.minL {
		return 0, nil
	}

	t := time.Now().Format(time.RFC3339)

	aux := struct {
		Level   string `json:"level"`
		Time    string `json:"time"`
		Context string `json:"context"`
		Message string `json:"message"`
		Trace   string `json:"trace,omitempty"`
	}{
		Level:   level.String(),
		Time:    t,
		Context: context,
		Message: message,
	}

	if level >= ERROR {
		aux.Trace = string(debug.Stack())
	}

	var line []byte

	line, err := json.Marshal(aux)

	if err != nil {
		message = fmt.Sprintf("%s %s", "unable to marshal log message:", err.Error())
		line = []byte(fmt.Sprintf("%s %v[%s]: %s", t, ERROR, "JsonLogger", message))
	}

	jl.mu.Lock()

	defer jl.mu.Unlock()

	return jl.o.Write(append(line, '\n'))
}
