// Package postgres handles database connection and initialization using GORM
// with PostgreSQL as the backend. It provides connection management, automatic
// schema migration, and a global database instance for backup operations.
package postgres

import (
	"fmt"
	"os"
	"time"

	"github.com/RubenPari/clear-songs/internal/infrastructure/persistence/postgres/models"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Db is the global GORM database instance, initialized by Init().
var Db *gorm.DB = nil

// Init initializes the database connection from environment variables (DB_HOST,
// DB_PORT, DB_USER, DB_PASSWORD, DB_NAME). Returns nil if configuration is missing
// or connection fails, allowing the application to continue without database backup.
func Init() error {
	// postgres credentials
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// Check if database configuration is provided
	if host == "" || port == "" || user == "" || password == "" || dbname == "" {
		zap.L().Warn("database environment variables not set, backup disabled")
		return nil // Return nil to allow application to continue without database
	}

	// create the connection string
	postgresInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", host, port, user, password, dbname)

	// Open the connection
	var errConnectDb error
	db, errConnectDb := gorm.Open(postgres.Open(postgresInfo), &gorm.Config{})

	if errConnectDb != nil {
		zap.L().Warn("database connection failed, backup disabled", zap.Error(errConnectDb))
		return nil // Return nil to allow application to continue without database
	}

	// Extract the underlying sql.DB to configure connection pooling
	sqlDB, err := db.DB()
	if err != nil {
		zap.L().Warn("failed to extract sql.DB for pooling", zap.Error(err))
	} else {
		// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
		sqlDB.SetMaxIdleConns(10)
		// SetMaxOpenConns sets the maximum number of open connections to the database.
		sqlDB.SetMaxOpenConns(100)
		// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
		sqlDB.SetConnMaxLifetime(time.Hour)
	}

	// test connection
	errTestDb := db.Exec("SELECT 1").Error

	if errTestDb != nil {
		zap.L().Warn("database connection test failed, backup disabled", zap.Error(errTestDb))
		return nil // Return nil to allow application to continue without database
	}

	// auto-migration
	errMigration := db.AutoMigrate(
		&models.TrackDB{},
	)

	if errMigration != nil {
		zap.L().Warn("database migration failed, backup disabled", zap.Error(errMigration))
		return nil // Return nil to allow application to continue without database
	}

	Db = db

	zap.L().Info("connected to database with pooling configured")

	return nil
}
