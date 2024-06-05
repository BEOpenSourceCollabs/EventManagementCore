package main

import (
	"net/http"

	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/config"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/net"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/net/routes"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/persist"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/repository"
)

func main() {
	envConfig := config.NewEnvironmentConfiguration()

	// initialize database from environment configuration
	database, err := persist.NewDatabase(envConfig.Database)
	if err != nil {
		panic(err)
	}

	// Create a new instance of the appRouter
	router := net.NewAppRouter()

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
