package middleware

import (
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/logger"
)

type RequestLoggerMiddleware struct{}

func (rmw RequestLoggerMiddleware) BeforeNext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wr := httptest.NewRecorder()
		next.ServeHTTP(wr, r)
		end := time.Since(start)

		logger.AppLogger.InfoF("RequestLoggerMiddleware", "Method: %s | Path: %s | Status: %d | Time: %v", r.Method, r.URL.Path, wr.Code, end)

		// Copy headers from recorder to response
		for k, vs := range wr.Header() {
			for _, v := range vs {
				w.Header().Add(k, v)
			}
		}

		// write recorders code and body to the response
		w.WriteHeader(wr.Code)
		w.Write(wr.Body.Bytes())
	})
}
