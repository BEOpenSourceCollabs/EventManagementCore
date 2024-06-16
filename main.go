package main

import (
	"fmt"
	"net/http"

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
	// Apply default middleware
	router.Use(
		middleware.RequestLoggerMiddleware{},
	)

	// Apply development middleware(s)
	if !envConfig.Env.IsProduction() {
		router.Use(
			middleware.CorsMiddleware{
				Options: middleware.WideOpen,
			},
		)
	}

	// Register static files handler
	router.Get("/", fs)

	// Register a GET route for the root URL ("/")
	router.Get("/api/health", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Ok"))
		},
	)) // Use the WideOpen CORS options to allow unrestricted access

	userRepo := repository.NewSQLUserRepository(
		database,
	)

	// initialize and mount routes
	routes.NewUserRoutes(
		router,
		userRepo,
		&envConfig.Security.Authentication,
	)

	authService := service.NewAuthService(&envConfig.Security.Authentication, userRepo)

	routes.NewAuthRoutes(
		router,
		authService,
		&envConfig.Security.Authentication,
	)

	address := fmt.Sprintf(":%d", envConfig.Port)
	logger.AppLogger.InfoF("main", "Starting server in %s mode on %s", envConfig.Env, address)
	logger.AppLogger.Fatal("main", http.ListenAndServe(address, router))
}
