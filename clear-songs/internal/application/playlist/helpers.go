package playlist

import (
	"context"
	"time"

	"github.com/RubenPari/clear-songs/internal/domain/shared"
	spotifyAPI "github.com/zmb3/spotify"
	"go.uber.org/zap"
)

// extractPlaylistTrackIDs extracts the track IDs from a slice of playlist tracks.
func extractPlaylistTrackIDs(tracks []spotifyAPI.PlaylistTrack) []spotifyAPI.ID {
	trackIDs := make([]spotifyAPI.ID, 0, len(tracks))
	for _, track := range tracks {
		trackIDs = append(trackIDs, track.Track.ID)
	}
	return trackIDs
}

// fetchPlaylistTracks returns the tracks for a playlist, reading from cache when available
// and falling back to the Spotify API. The result is written back to the cache.
func fetchPlaylistTracks(
	ctx context.Context,
	spotifyRepo shared.SpotifyRepository,
	cacheRepo shared.CacheRepository,
	playlistID spotifyAPI.ID,
) ([]spotifyAPI.PlaylistTrack, error) {
	cached, err := cacheRepo.GetPlaylistTracks(ctx, playlistID)
	if err == nil && len(cached) > 0 {
		return cached, nil
	}

	tracks, err := spotifyRepo.GetAllPlaylistTracks(ctx, playlistID)
	if err != nil {
		return nil, err
	}

	if err := cacheRepo.SetPlaylistTracks(ctx, playlistID, tracks, 5*time.Minute); err != nil {
		zap.L().Warn("failed to cache playlist tracks", zap.String("playlist_id", playlistID.String()), zap.Error(err))
	}
	return tracks, nil
}

const fetchPlaylistTracksCacheTTL = 5 * time.Minute
