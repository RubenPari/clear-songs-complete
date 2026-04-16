package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/RubenPari/clear-songs/internal/domain/shared/utils"
	"github.com/RubenPari/clear-songs/internal/infrastructure/di"
	"github.com/RubenPari/clear-songs/internal/infrastructure/persistence/postgres"
	httptransport "github.com/RubenPari/clear-songs/internal/infrastructure/transport/http"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// main starts the application.
func main() {
	// Initialize environment and DI
	utils.LoadEnvVariables()

	// Switch to release mode in production based on env
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Initialize Database with Pooling
	log.Println("Initializing database...")
	if errConnectDb := postgres.Init(); errConnectDb != nil {
		log.Printf("WARNING: Database initialization failed: %v", errConnectDb)
	}

	log.Println("Initializing DI container...")
	container, err := di.NewContainer()
	if err != nil {
		log.Fatalf("Failed to initialize DI container: %v", err)
	}

	// Setup Gin Router
	log.Println("Setting up router...")
	router := gin.Default()

	allowedOrigins := []string{"http://127.0.0.1", "http://127.0.0.1:4200", "http://localhost:4200"}
	if frontendURL := os.Getenv("FRONTEND_URL"); frontendURL != "" {
		allowedOrigins = append(allowedOrigins, frontendURL)
	}

	router.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Setup Routes
	httptransport.SetUpRoutes(router, container)

	// Create HTTP server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	// Long summaries (many sequential Gemini calls) can exceed 2m; default write timeout 6m.
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
		Addr:    ":" + port,
		Handler: router,
		ReadTimeout:  timeoutDur(readSec),
		WriteTimeout: timeoutDur(writeSec),
		IdleTimeout:  120 * time.Second,
	}
	log.Printf("HTTP server timeouts: read=%v write=%v (override with HTTP_READ_TIMEOUT_SEC / HTTP_WRITE_TIMEOUT_SEC)", srv.ReadTimeout, srv.WriteTimeout)

	// Run server in a goroutine so that it doesn't block
	go func() {
		log.Println("Starting server on :" + port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	// Properly close database connections
	if postgres.Db != nil {
		sqlDB, err := postgres.Db.DB()
		if err == nil {
			sqlDB.Close()
			log.Println("Database connection closed")
		}
	}

	log.Println("Server exiting")
}
