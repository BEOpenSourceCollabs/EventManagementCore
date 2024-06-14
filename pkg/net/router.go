// Package net provides network-related utilities
package net

import (
	"fmt"
	"net/http"

	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/logger"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/net/middleware"
)

// AppRouter is a custom router that embeds the http.ServeMux
type AppRouter struct {
	*http.ServeMux // ServeMux is an HTTP request multiplexer that matches the URL of each incoming request against a list of registered patterns
	middleware     []middleware.Middleware
}

// NewAppRouter creates a new AppRouter instance with an initialized ServeMux
func NewAppRouter() AppRouter {
	return AppRouter{
		ServeMux: http.NewServeMux(), // Initialize a new ServeMux
	}
}

func (r *AppRouter) handle(pattern string, handler http.Handler) {
	for _, middleware := range r.middleware {
		handler = middleware.BeforeNext(handler)
	}
	r.ServeMux.Handle(pattern, handler)
}

func (r *AppRouter) Use(middlewares ...middleware.Middleware) {
	r.middleware = append(r.middleware, middlewares...)
}

// Get registers a handler for GET requests with the given pattern
func (r *AppRouter) Get(pattern string, handler http.Handler) {
	logger.AppLogger.InfoF("AppRouter", "Mapped [GET] %s\n", pattern)
	r.handle(fmt.Sprintf("GET %s", pattern), handler)
}

// Post registers a handler for POST requests with the given pattern
func (r *AppRouter) Post(pattern string, handler http.Handler) {
	logger.AppLogger.InfoF("AppRouter", "Mapped [POST] %s\n", pattern)
	r.handle(fmt.Sprintf("POST %s", pattern), handler)
}

// Update registers a handler for UPDATE requests with the given pattern
func (r *AppRouter) Update(pattern string, handler http.Handler) {
	logger.AppLogger.InfoF("AppRouter", "Mapped [UPDATE] %s\n", pattern)
	r.handle(fmt.Sprintf("UPDATE %s", pattern), handler)
}

// Put registers a handler for PUT requests with the given pattern
func (r *AppRouter) Put(pattern string, handler http.Handler) {
	logger.AppLogger.InfoF("AppRouter", "Mapped [PUT] %s\n", pattern)
	r.handle(fmt.Sprintf("PUT %s", pattern), handler)
}

// Delete registers a handler for DELETE requests with the given pattern
func (r *AppRouter) Delete(pattern string, handler http.Handler) {
	logger.AppLogger.InfoF("AppRouter", "Mapped [DELETE] %s\n", pattern)
	r.handle(fmt.Sprintf("DELETE %s", pattern), handler)
}

// Options registers a handler for OPTIONS requests with the given pattern
func (r *AppRouter) Options(pattern string, handler http.Handler) {
	logger.AppLogger.InfoF("AppRouter", "Mapped [OPTIONS] %s\n", pattern)
	r.handle(fmt.Sprintf("OPTIONS %s", pattern), handler)
}
