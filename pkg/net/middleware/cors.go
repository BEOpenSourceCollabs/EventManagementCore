// Package middleware contains HTTP middleware utilities
package middleware

import (
	"fmt"
	"net/http"
	"strings"
)

// CorsOptions defines the structure for CORS configuration options
type CorsOptions struct {
	Origin           string   // Origin specifies the origin that is allowed to access the resource
	Methods          []string // Methods is a list of HTTP methods allowed when accessing the resource
	Headers          []string // Headers is a list of HTTP headers that can be used when making the actual request
	ExposeHeaders    []string // ExposeHeaders indicates which headers are safe to expose to the API of a CORS API specification
	AllowCredentials bool     // AllowCredentials indicates whether or not the response to the request can be exposed when the credentials flag is true
}

// WideOpen is a preset configuration allowing unrestricted CORS access
var WideOpen = CorsOptions{
	Origin:           "*",                                                                                      // Allow any origin
	Methods:          []string{"POST", "GET", "OPTIONS", "PUT", "UPDATE", "DELETE"},                            // Allow common HTTP methods
	Headers:          []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "Authorization"}, // Allow common HTTP headers
	ExposeHeaders:    []string{"*"},                                                                            // Allow any headers to be exposed
	AllowCredentials: true,                                                                                     // Allow credentials to be included
}

type CorsMiddleware struct {
	Options CorsOptions
}

func (cmw CorsMiddleware) BeforeNext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", cmw.Options.Origin)
		w.Header().Set("Access-Control-Allow-Methods", strings.Join(cmw.Options.Methods, ", "))
		w.Header().Set("Access-Control-Allow-Headers", strings.Join(cmw.Options.Headers, ", "))
		w.Header().Set("Access-Control-Expose-Headers", fmt.Sprint(cmw.Options.ExposeHeaders))
		w.Header().Set("Access-Control-Allow-Credentials", fmt.Sprint(cmw.Options.AllowCredentials))
		next.ServeHTTP(w, r)
	})
}

// CorsMiddleware is an HTTP middleware function that adds CORS headers to the response
// func CorsMiddleware(next http.Handler, options CorsOptions) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		logger.AppLogger.InfoF("CorsMiddleware", "adding cors headers with options - %v", options)
// 		w.Header().Set("Access-Control-Allow-Origin", options.Origin)
// 		w.Header().Set("Access-Control-Allow-Methods", strings.Join(options.Methods, ", "))
// 		w.Header().Set("Access-Control-Allow-Headers", strings.Join(options.Headers, ", "))
// 		w.Header().Set("Access-Control-Expose-Headers", fmt.Sprint(options.ExposeHeaders))
// 		w.Header().Set("Access-Control-Allow-Credentials", fmt.Sprint(options.AllowCredentials))
// 		next.ServeHTTP(w, r)
// 	})
// }
