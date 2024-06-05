// Package net provides network-related utilities
package net

import (
	"fmt"
	"net/http"
)

// AppRouter is a custom router that embeds the http.ServeMux
type AppRouter struct {
	*http.ServeMux // ServeMux is an HTTP request multiplexer that matches the URL of each incoming request against a list of registered patterns
}

// NewAppRouter creates a new AppRouter instance with an initialized ServeMux
func NewAppRouter() AppRouter {
	return AppRouter{
		ServeMux: http.NewServeMux(), // Initialize a new ServeMux
	}
}

// Get registers a handler for GET requests with the given pattern
func (r *AppRouter) Get(pattern string, handler http.Handler) {
	fmt.Printf("[GET] %s\n", pattern)
	r.ServeMux.Handle(fmt.Sprintf("GET %s", pattern), handler)
}

// Post registers a handler for POST requests with the given pattern
func (r *AppRouter) Post(pattern string, handler http.Handler) {
	fmt.Printf("[POST] %s\n", pattern)
	r.ServeMux.Handle(fmt.Sprintf("POST %s", pattern), handler)
}

// Update registers a handler for UPDATE requests with the given pattern
func (r *AppRouter) Update(pattern string, handler http.Handler) {
	fmt.Printf("[UPDATE] %s\n", pattern)
	r.ServeMux.Handle(fmt.Sprintf("UPDATE %s", pattern), handler)
}

// Put registers a handler for PUT requests with the given pattern
func (r *AppRouter) Put(pattern string, handler http.Handler) {
	fmt.Printf("[PUT] %s\n", pattern)
	r.ServeMux.Handle(fmt.Sprintf("PUT %s", pattern), handler)
}

// Delete registers a handler for DELETE requests with the given pattern
func (r *AppRouter) Delete(pattern string, handler http.Handler) {
	fmt.Printf("[DELETE] %s\n", pattern)
	r.ServeMux.Handle(fmt.Sprintf("DELETE %s", pattern), handler)
}

// Options registers a handler for OPTIONS requests with the given pattern
func (r *AppRouter) Options(pattern string, handler http.Handler) {
	fmt.Printf("[OPTIONS] %s\n", pattern)
	r.ServeMux.Handle(fmt.Sprintf("OPTIONS %s", pattern), handler)
}
