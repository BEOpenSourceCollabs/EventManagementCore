package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/config"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/net"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/net/middleware"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/net/routes"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/persist"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/repository"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/service"
)

func main() {
	envConfig := config.NewEnvironmentConfiguration()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// Create a static handler to serve contents of 'static' folder.
	fs := http.FileServer(http.Dir("./static"))

	// initialize database from environment configuration
	database, err := persist.NewDatabase(envConfig.Database)
	if err != nil {
		logger.Fatal(err)
	}

	// Create a new instance of the appRouter
	router := net.NewAppRouter()

	// Register static files handler
	router.Get("/", fs)

	// Register a GET route for the root URL ("/") with CORS middleware
	router.Get("/api/health", middleware.CorsMiddleware(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Ok"))
		},
	), middleware.WideOpen)) // Use the WideOpen CORS options to allow unrestricted access

	userRepo := repository.NewSQLUserRepository(
		database,
	)

	// initialize and mount routes
	routes.NewUserRoutes(
		router,
		userRepo,
	)

	authService := service.NewAuthService(&envConfig, logger, userRepo)
	routes.NewAuthRoutes(router, authService, logger)

	// http server configured with some defaults
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", envConfig.Port),
		Handler:      middleware.RequestLoggerMiddleware(router, logger),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Printf("Starting server in %s mode on %s", envConfig.Env, srv.Addr)

	// Start the HTTP server using the router
	err = srv.ListenAndServe()

	logger.Fatal(err)

}
