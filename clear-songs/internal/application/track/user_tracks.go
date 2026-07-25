package track

import (
	"context"
	"fmt"
	"time"

	"github.com/RubenPari/clear-songs/internal/domain/shared"
	spotifyAPI "github.com/zmb3/spotify"
	"go.uber.org/zap"
)

// userTracksCacheTTL is the duration for which the user's saved tracks are cached
// in Redis before being re-fetched from the Spotify API.
const userTracksCacheTTL = 5 * time.Minute

// getUserTracks retrieves the user's saved tracks, first checking the cache and
// falling back to the Spotify API on a cache miss. The result is written back to
// the cache with a 5-minute TTL.
func getUserTracks(ctx context.Context, spotifyRepo shared.SpotifyRepository, cacheRepo shared.CacheRepository) ([]spotifyAPI.SavedTrack, error) {
	cached, err := cacheRepo.GetUserTracks(ctx)
	if err == nil && cached != nil {
		return cached, nil
	}
	if err != nil {
		zap.L().Warn("failed to read user tracks cache", zap.Error(err))
	}

	tracks, err := spotifyRepo.GetAllUserTracks(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetch user tracks from spotify: %w", err)
	}

	if err := cacheRepo.SetUserTracks(ctx, tracks, userTracksCacheTTL); err != nil {
		zap.L().Warn("failed to write user tracks cache", zap.Error(err))
	}

	return tracks, nil
}
