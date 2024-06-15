package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/config"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/logger"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/net"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/net/middleware"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/net/routes"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/persist"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/repository"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/service"
)

func main() {
	envConfig := config.NewEnvironmentConfiguration()

	// Create a static handler to serve contents of 'static' folder.
	fs := http.FileServer(http.Dir("./static"))

	// initialize database from environment configuration
	database, err := persist.NewDatabase(envConfig.Database)
	if err != nil {
		logger.AppLogger.Fatal("main", err)
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

	eventRepo := repository.NewSQLEventRepository(
		database,
	)

	// initialize and mount routes
	routes.NewEventRoutes(
		router,
		eventRepo,
	)

	authService := service.NewAuthService(&envConfig.Security.Authentication, userRepo)

	routes.NewAuthRoutes(
		router,
		authService,
		&envConfig.Security.Authentication,
	)

	// http server configured with some defaults
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", envConfig.Port),
		Handler:      middleware.RequestLoggerMiddleware(router),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.AppLogger.InfoF("main", "Starting server in %s mode on %s", envConfig.Env, srv.Addr)

	// Start the HTTP server using the router
	err = srv.ListenAndServe()

	logger.AppLogger.Fatal("main", err)

}
