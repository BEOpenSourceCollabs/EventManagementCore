package middleware

import (
	"log"
	"net/http"
	"time"
)

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}

func NewStatusRecorder(w http.ResponseWriter) *statusRecorder {
	return &statusRecorder{w, http.StatusOK}
}

func RequestLoggerMiddleware(next http.Handler, logger *log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// create statusRecorder to capture status code
		sr := NewStatusRecorder(w)

		next.ServeHTTP(sr, r)

		end := time.Since(start)

		logger.Printf("Path: %s | Status: %d | Time: %v", r.URL.Path, sr.statusCode, end)
	})
}
