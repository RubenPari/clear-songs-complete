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

	"github.com/RubenPari/clear-songs/internal/domain/shared/utils"
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

// main starts the application.
func main() {
	logger := logging.InitFromEnv()
	defer logging.SafeSync(logger)

	utils.LoadEnvVariables()

	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	logger.Info("initializing database")
	if errConnectDb := postgres.Init(); errConnectDb != nil {
		logger.Warn("database initialization failed", zap.Error(errConnectDb))
	}

	logger.Info("initializing DI container")
	container, err := di.NewContainer()
	if err != nil {
		logger.Fatal("failed to initialize DI container", zap.Error(err))
	}

	logger.Info("setting up router")
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

	allowedOrigins := []string{"http://127.0.0.1", "http://127.0.0.1:4200", "http://localhost:4200"}
	if frontendURL := os.Getenv("FRONTEND_URL"); frontendURL != "" {
		allowedOrigins = append(allowedOrigins, frontendURL)
	}

	router.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	httptransport.SetUpRoutes(router, container)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	writeSec := 360
	if s := strings.TrimSpace(os.Getenv("HTTP_WRITE_TIMEOUT_SEC")); s != "" {
		if n, err := strconv.Atoi(s); err == nil && n >= 30 && n <= 3600 {
			writeSec = n
		}
	}
	readSec := writeSec
	if s := strings.TrimSpace(os.Getenv("HTTP_READ_TIMEOUT_SEC")); s != "" {
		if n, err := strconv.Atoi(s); err == nil && n >= 30 && n <= 3600 {
			readSec = n
		}
	}
	timeoutDur := func(sec int) time.Duration { return time.Duration(sec) * time.Second }

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  timeoutDur(readSec),
		WriteTimeout: timeoutDur(writeSec),
		IdleTimeout:  120 * time.Second,
	}
	logger.Info("HTTP server timeouts configured",
		zap.Duration("read_timeout", srv.ReadTimeout),
		zap.Duration("write_timeout", srv.WriteTimeout),
	)

	go func() {
		logger.Info("starting server", zap.String("port", port))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("server error", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("server forced to shutdown", zap.Error(err))
	}

	if postgres.Db != nil {
		sqlDB, err := postgres.Db.DB()
		if err == nil {
			_ = sqlDB.Close()
			logger.Info("database connection closed")
		}
	}

	logger.Info("server exiting")
}
