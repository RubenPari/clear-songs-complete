package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/RubenPari/clear-songs/internal/infrastructure/config"
	"github.com/RubenPari/clear-songs/internal/infrastructure/di"
	"github.com/RubenPari/clear-songs/internal/infrastructure/logging"
	"github.com/RubenPari/clear-songs/internal/infrastructure/persistence/postgres"
	httptransport "github.com/RubenPari/clear-songs/internal/infrastructure/transport/http"
	"github.com/RubenPari/clear-songs/internal/infrastructure/transport/http/middleware"
	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Server configuration default values and limits.
const (
	defaultPort           = "3000"
	defaultTimeoutSeconds = 360
	minTimeoutSeconds     = 30
	maxTimeoutSeconds     = 3600
	idleTimeoutSeconds    = 120
	gracefulShutdownSec   = 5
)

// Default origins permitted for local development.
var (
	allowedOriginsLocal = []string{
		"http://127.0.0.1",
		"http://127.0.0.1:4200",
		"http://localhost:4200",
	}
)

// parseTimeout reads an integer duration from an environment variable in seconds.
// It falls back to defaultSec if the variable is missing, invalid, or out of bounds.
func parseTimeout(envVar string, defaultSec int) int {
	s := strings.TrimSpace(os.Getenv(envVar))
	if s == "" {
		return defaultSec
	}
	n, err := strconv.Atoi(s)
	if err != nil || n < minTimeoutSeconds || n > maxTimeoutSeconds {
		return defaultSec
	}
	return n
}

// setupCORSConfig builds the CORS configuration including local development URLs
// and an optional frontend URL from the environment configuration.
func setupCORSConfig(logger *zap.Logger) cors.Config {
	origins := allowedOriginsLocal
	if frontendURL := config.GetFrontendURL(); frontendURL != "" {
		origins = append(origins, frontendURL)
		logger.Debug("added custom frontend URL to CORS origins", zap.String("url", frontendURL))
	}

	return cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
}

// setupRouter initializes the Gin engine with logging, recovery, request ID, and CORS middleware.
func setupRouter(logger *zap.Logger) *gin.Engine {
	router := gin.New()
	router.Use(middleware.RequestIDMiddleware())
	router.Use(ginzap.GinzapWithConfig(logger, &ginzap.Config{
		TimeFormat: time.RFC3339,
		UTC:        true,
		Context: ginzap.Fn(func(c *gin.Context) []zapcore.Field {
			return []zapcore.Field{zap.String("request_id", c.GetString(logging.RequestIDKey))}
		}),
	}))
	router.Use(ginzap.RecoveryWithZap(logger, true))
	router.Use(cors.New(setupCORSConfig(logger)))

	return router
}

// createHTTPServer instantiates an http.Server configured with custom timeouts.
func createHTTPServer(router *gin.Engine, port string, readTimeoutSec, writeTimeoutSec int) *http.Server {
	return &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  time.Duration(readTimeoutSec) * time.Second,
		WriteTimeout: time.Duration(writeTimeoutSec) * time.Second,
		IdleTimeout:  time.Duration(idleTimeoutSeconds) * time.Second,
	}
}

// closeDatabase safely closes the underlying database connection pool if active.
func closeDatabase(logger *zap.Logger) {
	if postgres.Db != nil {
		sqlDB, err := postgres.Db.DB()
		if err == nil {
			_ = sqlDB.Close()
			logger.Info("database connection closed")
		}
	}
}

func main() {
	// Initialize structured logging and flush buffer on exit.
	logger := logging.InitFromEnv()
	defer logging.SafeSync(logger)

	// Load environment variables and configure framework modes.
	config.LoadEnvVariables()
	setGinMode()

	// Initialize database connection (optional backup store).
	if err := postgres.Init(); err != nil {
		logger.Warn("database initialization failed", zap.Error(err))
	}
	defer closeDatabase(logger)

	// Initialize dependency injection container (requires Redis).
	container, err := di.NewContainer()
	if err != nil {
		logger.Fatal("failed to initialize DI container", zap.Error(err))
	}

	// Read port and timeout configuration from environment variables.
	port := getPort()
	readTimeoutSec := parseTimeout("HTTP_READ_TIMEOUT_SEC", defaultTimeoutSeconds)
	writeTimeoutSec := parseTimeout("HTTP_WRITE_TIMEOUT_SEC", defaultTimeoutSeconds)

	// Set up router, middleware, and application HTTP routes.
	router := setupRouter(logger)
	httptransport.SetUpRoutes(router, container)

	// Create and configure the HTTP server.
	srv := createHTTPServer(router, port, readTimeoutSec, writeTimeoutSec)
	logger.Info("HTTP server configured",
		zap.String("port", port),
		zap.Duration("read_timeout", srv.ReadTimeout),
		zap.Duration("write_timeout", srv.WriteTimeout),
	)

	// Start the server asynchronously and await termination signals.
	startServer(logger, srv)
	awaitShutdown(logger, srv)
}

// setGinMode sets Gin execution mode depending on GIN_MODE env var.
func setGinMode() {
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
}

// getPort returns the PORT env variable or defaults to defaultPort.
func getPort() string {
	if port := os.Getenv("PORT"); port != "" {
		return port
	}
	return defaultPort
}

// startServer runs the HTTP server in a separate goroutine.
func startServer(logger *zap.Logger, srv *http.Server) {
	go func() {
		logger.Info("starting server", zap.String("port", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("server error", zap.Error(err))
		}
	}()
}

// awaitShutdown blocks until an OS termination signal is received,
// then triggers a graceful shutdown of the HTTP server.
func awaitShutdown(logger *zap.Logger, srv *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), gracefulShutdownSec*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("server forced to shutdown", zap.Error(err))
	}

	logger.Info("server exited gracefully")
}
