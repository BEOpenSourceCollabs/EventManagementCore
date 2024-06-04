package main

import (
	"net/http"

	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/config"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/net"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/net/middleware"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/persist"
)

func main() {
	envConfig := config.NewEnvironmentConfiguration()

	/* todo: database */
	_, err := persist.NewDatabase(envConfig.Database)
	if err != nil {
		panic(err)
	}

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
