package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/config"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/logging"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/net"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/net/middleware"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/net/routes"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/persist"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/repository"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/service"
)

func main() {
	envConfig := config.NewEnvironmentConfiguration()

	// initialize log writer that outs logs as either json or text to stdout
	var lw logging.LogWriter
	if envConfig.Env.IsProduction() {
		lw = logging.NewJsonLogWriter(os.Stdout, logging.INFO)
	} else {
		lw = logging.NewTextLogWriter(os.Stdout, logging.INFO)
	}
	// ContextLogger for main funcition, uses default log writer to print logs with Context
	mainLogger := logging.NewContextLogger(lw, "Main")

	mainLogger.Infof("loaded environment configuration for %s", envConfig.Env)

	// Create a static handler to serve contents of 'static' folder.
	fs := http.FileServer(http.Dir("./static"))

	// initialize database from environment configuration
	database, err := persist.NewDatabase(envConfig.Database)
	if err != nil {
		mainLogger.Fatal(err, "database connection error")
	}

	// Create a new instance of the appRouter
	router := net.NewAppRouter(lw)
	// Apply default middleware
	router.Use(
		middleware.RequestLoggerMiddleware{
			Logger: logging.NewContextLogger(lw, "RequestLoggerMiddleware"),
		},
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

	jwtService := service.NewJsonWebTokenService(
		&envConfig.Security.JsonWebToken,
		lw,
	)

	authService := service.NewJsonWebTokenAuthenticationService(
		userRepo,
		jwtService,
		lw,
		&service.AuthenticationServiceConfiguration{
			IsProduction: envConfig.Env.IsProduction(),
		},
	)

	// initialize and mount routes
	routes.NewJsonWebTokenUserRoutes(
		router,
		userRepo,
		&jwtService,
		lw,
	)

	routes.NewJsonWebTokenAuthenticationRoutes(
		router,
		authService,
		&jwtService,
		lw,
	)

	routes.NewGoogleAuthenticationRoutes(
		router,
		service.NewGoogleAuthenticationService(
			&envConfig.Security.Google,
			userRepo,
			jwtService,
			lw,
		),
		authService,
		lw,
	)

	address := fmt.Sprintf(":%d", envConfig.Port)
	mainLogger.Infof("Starting server in %s mode on %s", envConfig.Env, address)

	err = http.ListenAndServe(address, router)

	if err != nil {
		mainLogger.Fatal(err, "failed to start http server")
	}
}
