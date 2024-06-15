package middleware

import (
	"net/http"
	"time"

	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/logger"
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

func RequestLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// create statusRecorder to capture status code
		sr := NewStatusRecorder(w)

		next.ServeHTTP(sr, r)

		end := time.Since(start)

		logger.AppLogger.InfoF("RequestLoggerMiddleware", "Path: %s | Status: %d | Time: %v", r.URL.Path, sr.statusCode, end)
	})
}
