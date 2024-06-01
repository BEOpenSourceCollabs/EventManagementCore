// Package net provides network-related utilities
package net

import (
	"fmt"
	"net/http"
)

// appRouter is a custom router that embeds the http.ServeMux
type appRouter struct {
	*http.ServeMux // ServeMux is an HTTP request multiplexer that matches the URL of each incoming request against a list of registered patterns
}

// NewAppRouter creates a new appRouter instance with an initialized ServeMux
func NewAppRouter() appRouter {
	return appRouter{
		ServeMux: http.NewServeMux(), // Initialize a new ServeMux
	}
}

// Get registers a handler for GET requests with the given pattern
func (r *appRouter) Get(pattern string, handler http.Handler) {
	fmt.Printf("[GET] %s\n", pattern)
	r.ServeMux.Handle(fmt.Sprintf("GET %s", pattern), handler)
}

// Post registers a handler for POST requests with the given pattern
func (r *appRouter) Post(pattern string, handler http.Handler) {
	fmt.Printf("[POST] %s\n", pattern)
	r.ServeMux.Handle(fmt.Sprintf("POST %s", pattern), handler)
}

// Update registers a handler for UPDATE requests with the given pattern
func (r *appRouter) Update(pattern string, handler http.Handler) {
	fmt.Printf("[UPDATE] %s\n", pattern)
	r.ServeMux.Handle(fmt.Sprintf("UPDATE %s", pattern), handler)
}

// Put registers a handler for PUT requests with the given pattern
func (r *appRouter) Put(pattern string, handler http.Handler) {
	fmt.Printf("[PUT] %s\n", pattern)
	r.ServeMux.Handle(fmt.Sprintf("PUT %s", pattern), handler)
}

// Delete registers a handler for DELETE requests with the given pattern
func (r *appRouter) Delete(pattern string, handler http.Handler) {
	fmt.Printf("[DELETE] %s\n", pattern)
	r.ServeMux.Handle(fmt.Sprintf("DELETE %s", pattern), handler)
}

// Options registers a handler for OPTIONS requests with the given pattern
func (r *appRouter) Options(pattern string, handler http.Handler) {
	fmt.Printf("[OPTIONS] %s\n", pattern)
	r.ServeMux.Handle(fmt.Sprintf("OPTIONS %s", pattern), handler)
}
