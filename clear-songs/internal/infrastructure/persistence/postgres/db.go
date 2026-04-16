/**
 * Database Package
 *
 * This package handles database connection and initialization using GORM
 * (Go Object-Relational Mapping) with PostgreSQL as the database backend.
 *
 * The package provides:
 * - Database connection management
 * - Automatic schema migration
 * - Global database instance for use throughout the application
 *
 * Database Schema:
 * The package uses GORM's AutoMigrate feature to automatically create and
 * update database tables based on model definitions. Currently manages:
 * - TrackDB model: Stores track metadata and user library information
 *
 * Connection Configuration:
 * Database credentials are loaded from environment variables:
 * - DB_HOST: Database host address
 * - DB_PORT: Database port number
 * - DB_USER: Database username
 * - DB_PASSWORD: Database password
 * - DB_NAME: Database name
 *
 * @package postgres
 * @author Clear Songs Development Team
 */
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

/**
 * Global Database Instance
 *
 * This variable holds the GORM database connection instance.
 * It is initialized by the Init() function and can be accessed
 * throughout the application for database operations.
 *
 * The variable is set to nil initially and will be assigned
 * a valid database connection after successful initialization.
 */
var Db *gorm.DB = nil

// Initializes.
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
