package cmd

import (
	"flag"
	"os"
	"os/signal"
	"product-service/internal/adapter"
	"product-service/internal/infrastructure"
	"product-service/internal/route"
	"product-service/pkg/validator"
	"runtime"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// RunServer function is used to run the server
// It will initialize the logger, fiber app, and the routes
// It will also handle graceful shutdown
func RunServer(cmd *flag.FlagSet, args []string) {
	envs := infrastructure.Envs
	logLevel, err := zerolog.ParseLevel(envs.App.LogLevel)
	if err != nil {
		logLevel = zerolog.InfoLevel
	}

	infrastructure.InitializeLogger(envs.App.Environtment, "app.log", logLevel)

	var (
		flagAppPort = cmd.String("port", "3000", "Application port") // ex: go run main.go server -port=8080
	)

	if err := cmd.Parse(args); err != nil {
		log.Fatal().Err(err).Msg("Error while parsing flags")
	}

	app := fiber.New()

	// Application Middlewares

	// Rate limiter middleware
	if envs.App.Environtment == "production" {
		app.Use(limiter.New(limiter.Config{
			Max:        50,
			Expiration: 30 * time.Second,
		}))
	}

	// CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,PATCH,OPTIONS,HEAD",
		AllowHeaders: "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin,Authorization",
	}))
	// End Application Middlewares

	adapter.Adapters.Sync(
		adapter.WithRestServer(app),
		adapter.WithShopeefunProductPostgres(),
		adapter.WithValidator(validator.NewValidator()),
	)

	route.SetupRoutes(app)

	// Run server in goroutine
	go func() {
		if err := app.Listen(":" + *flagAppPort); err != nil {
			log.Fatal().Msgf("Error while starting server: %v", err)
		}
	}()
	// End Run server in goroutine

	// Handle graceful shutdown
	quit := make(chan os.Signal, 1)

	shutdownSignals := []os.Signal{os.Interrupt, syscall.SIGTERM, syscall.SIGINT}
	if runtime.GOOS == "windows" {
		shutdownSignals = []os.Signal{os.Interrupt}
	}

	signal.Notify(quit, shutdownSignals...)
	<-quit
	log.Info().Msg("Server is shutting down ...")

	err = adapter.Adapters.Unsync()
	if err != nil {
		log.Error().Msgf("Error while closing adapters: %v", err)
	}

	log.Info().Msg("Server gracefully stopped")
}
