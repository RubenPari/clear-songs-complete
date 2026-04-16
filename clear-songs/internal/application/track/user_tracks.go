package track

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/RubenPari/clear-songs/internal/domain/shared"
	spotifyAPI "github.com/zmb3/spotify"
)

const userTracksCacheTTL = 5 * time.Minute

// Fetches user tracks.
func getUserTracks(ctx context.Context, spotifyRepo shared.SpotifyRepository, cacheRepo shared.CacheRepository) ([]spotifyAPI.SavedTrack, error) {
	if cacheRepo != nil {
		cached, err := cacheRepo.GetUserTracks(ctx)
		if err == nil && cached != nil {
			return cached, nil
		}
		if err != nil {
			log.Printf("warning: failed to read user tracks cache: %v", err)
		}
	}

	tracks, err := spotifyRepo.GetAllUserTracks(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetch user tracks from spotify: %w", err)
	}

	if cacheRepo != nil {
		if err := cacheRepo.SetUserTracks(ctx, tracks, userTracksCacheTTL); err != nil {
			log.Printf("warning: failed to write user tracks cache: %v", err)
		}
	}

	return tracks, nil
}
