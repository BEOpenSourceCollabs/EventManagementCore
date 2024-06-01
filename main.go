package main

import (
	"net/http"

	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/net"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/net/middleware"
)

func main() {
	// Create a new instance of the appRouter
	router := net.NewAppRouter()

	// Register a GET route for the root URL ("/") with CORS middleware
	router.Get("/", middleware.CorsMiddleware(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello, world!"))
		},
	), middleware.WideOpen)) // Use the WideOpen CORS options to allow unrestricted access

	// Start the HTTP server on port 8081 using the router
	http.ListenAndServe(":8081", router)
}
