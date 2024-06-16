package middleware

import "net/http"

type Middleware interface {
	BeforeNext(next http.Handler) http.Handler
}
