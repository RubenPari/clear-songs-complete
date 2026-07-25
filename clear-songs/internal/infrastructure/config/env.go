// Package config provides environment variable loading and configuration helpers.
// It reads .env files and exposes functions for retrieving OAuth redirect URLs
// and frontend URLs.
package config

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

// GetRedirectURL returns the OAuth redirect URL from the environment.
// It accepts either REDIRECT_URL or REDIRECT_URI for backward compatibility.
func GetRedirectURL() string {
	redirectURL := os.Getenv("REDIRECT_URL")
	if redirectURL == "" {
		redirectURL = os.Getenv("REDIRECT_URI")
	}
	return redirectURL
}

// GetFrontendURL returns the frontend URL from the environment.
// It falls back to the local development URL when FRONTEND_URL is not set.
func GetFrontendURL() string {
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://127.0.0.1:4200"
	}
	return frontendURL
}

// LoadEnvVariables loads .env from the working directory (or its parent on Windows)
// and validates that a redirect URL is configured.
func LoadEnvVariables() {
	cwd, errCwd := os.Getwd()
	if errCwd != nil {
		zap.L().Fatal("error getting current working directory", zap.Error(errCwd))
	}

	if runtime.GOOS == "windows" {
		cwd = filepath.Dir(cwd)
	}

	envPath := filepath.Join(cwd, ".env")

	zap.L().Info("loading environment variables", zap.String("path", envPath))

	if errLoadFilePath := godotenv.Load(envPath); errLoadFilePath != nil {
		zap.L().Warn("error loading .env file, using system environment variables",
			zap.String("path", envPath),
			zap.Error(errLoadFilePath),
		)
	}

	zap.L().Info("loaded environment variables")

	redirectURL := GetRedirectURL()
	if redirectURL == "" {
		zap.L().Fatal("REDIRECT_URL or REDIRECT_URI not found in environment variables after loading .env file")
	}
	zap.L().Info("OAuth redirect URL configured", zap.String("redirect_url", redirectURL))
}
