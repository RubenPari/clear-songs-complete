package utils

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
	spotifyAPI "github.com/zmb3/spotify"
	"go.uber.org/zap"
)

// Converts tracks to id.
func ConvertTracksToID(tracks interface{}) ([]spotifyAPI.ID, error) {
	var trackIDs []spotifyAPI.ID

	switch t := tracks.(type) {
	case []spotifyAPI.FullTrack:
		for _, track := range t {
			trackIDs = append(trackIDs, track.ID)
		}
	case []spotifyAPI.PlaylistTrack:
		for _, track := range t {
			trackIDs = append(trackIDs, track.Track.ID)
		}
	case []spotifyAPI.SavedTrack:
		for _, track := range t {
			trackIDs = append(trackIDs, track.FullTrack.ID)
		}
	case []spotifyAPI.SavedAlbum:
		for _, album := range t {
			for _, track := range album.Tracks.Tracks {
				trackIDs = append(trackIDs, track.ID)
			}
		}
	default:
		return nil, errors.New(" ConvertTracksToID: Type input not supported")
	}

	return trackIDs, nil
}

// Loads env variables.
func LoadEnvVariables() {
	// get current working directory
	cwd, errCwd := os.Getwd()

	if errCwd != nil {
		zap.L().Fatal("error getting current working directory", zap.Error(errCwd))
	}

	// check if the OS is Windows
	if runtime.GOOS == "windows" {
		// move up one level folder
		cwd = filepath.Dir(cwd)
	}

	envPath := filepath.Join(cwd, ".env")

	zap.L().Info("loading environment variables", zap.String("path", envPath))

	errLoadFilePath := godotenv.Load(envPath)

	if errLoadFilePath != nil {
		zap.L().Warn("error loading .env file, using system environment variables",
			zap.String("path", envPath),
			zap.Error(errLoadFilePath),
		)
	}

	zap.L().Info("loaded environment variables")

	// Verify critical environment variables are loaded
	redirectURL := os.Getenv("REDIRECT_URL")
	if redirectURL == "" {
		redirectURL = os.Getenv("REDIRECT_URI")
	}
	if redirectURL == "" {
		zap.L().Fatal("REDIRECT_URL or REDIRECT_URI not found in environment variables after loading .env file")
	}
	zap.L().Info("OAuth redirect URL configured", zap.String("redirect_url", redirectURL))
}
