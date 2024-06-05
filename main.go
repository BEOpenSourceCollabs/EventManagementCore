package main

import (
	"net/http"

	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/config"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/net"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/net/middleware"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/net/routes"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/persist"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/repository"
)

func main() {
	envConfig := config.NewEnvironmentConfiguration()

	// Create a static handler to serve contents of 'static' folder.
	fs := http.FileServer(http.Dir("./static"))

	// initialize database from environment configuration
	database, err := persist.NewDatabase(envConfig.Database)
	if err != nil {
		panic(err)
	}

	// Create a new instance of the appRouter
	router := net.NewAppRouter()

	// Register static files handler
	router.Get("/", fs)

	// Register a GET route for the root URL ("/") with CORS middleware
	router.Get("/health", middleware.CorsMiddleware(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Ok"))
		},
	), middleware.WideOpen)) // Use the WideOpen CORS options to allow unrestricted access

	// initialize and mount routes
	routes.NewUserRoutes(
		router,
		repository.NewUserRepository(
			database,
		),
	)

	// initialize and mount routes
	routes.NewUserRoutes(
		router,
		repository.NewUserRepository(
			database,
		),
	)

	// Start the HTTP server on port 8081 using the router
	http.ListenAndServe(":8081", router)
}
