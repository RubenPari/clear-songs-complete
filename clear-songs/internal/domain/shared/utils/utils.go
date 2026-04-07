package utils

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/RubenPari/clear-songs/internal/infrastructure/persistence/postgres/models"
	"github.com/joho/godotenv"
	spotifyAPI "github.com/zmb3/spotify"
)

// ConvertTracksToID converts a list of tracks
// can be of type:
// - []spotifyAPI.FullTrack,
// - []spotifyAPI.PlaylistTrack,
// - []spotifyAPI.SavedTrack,
// - []spotifyAPI.SavedAlbum
// to a list of track IDs
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

// LoadEnvVariables loads the environment variables from the .env file in the
// current working directory.
func LoadEnvVariables() {
	// get current working directory
	cwd, errCwd := os.Getwd()

	if errCwd != nil {
		log.Fatalf("error getting current working directory: %v", errCwd)
	}

	// check if the OS is Windows
	if runtime.GOOS == "windows" {
		// move up one level folder
		cwd = filepath.Dir(cwd)
	}

	envPath := filepath.Join(cwd, ".env")

	log.Printf("Loading environment variables from: %s", envPath)

	errLoadFilePath := godotenv.Load(envPath)

	if errLoadFilePath != nil {
		log.Printf("Warning: error loading .env file from %s: %v. Using system environment variables.", envPath, errLoadFilePath)
	}

	log.Println("Loaded environment variables from .env file or system")

	// Verify critical environment variables are loaded
	redirectURL := os.Getenv("REDIRECT_URL")
	if redirectURL == "" {
		redirectURL = os.Getenv("REDIRECT_URI")
	}
	if redirectURL == "" {
		log.Fatal("REDIRECT_URL or REDIRECT_URI not found in environment variables after loading .env file")
	}
	log.Printf("OAuth Redirect URL configured: %s", redirectURL)
}

// FilterSummaryByRange is kept here for now as it's purely logic, 
// though it could be moved to application layer.
func FilterSummaryByRange(tracks []models.ArtistSummary, min int, max int) []models.ArtistSummary {
	var newTracks []models.ArtistSummary

	for _, track := range tracks {
		if min == 0 && max == 0 {
			newTracks = append(newTracks, track)
		} else if min == 0 && track.Count <= max {
			newTracks = append(newTracks, track)
		} else if max == 0 && track.Count >= min {
			newTracks = append(newTracks, track)
		} else if track.Count >= min && track.Count <= max {
			newTracks = append(newTracks, track)
		}
	}

	return newTracks
}
