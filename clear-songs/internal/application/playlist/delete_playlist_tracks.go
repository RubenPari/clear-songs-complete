// Package playlist implements application-layer use cases for playlist management.
// It orchestrates Spotify API calls and caching strategies to provide playlist
// retrieval, track deletion, and combined playlist/library cleanup operations.
package playlist

import (
	"context"

	"github.com/RubenPari/clear-songs/internal/domain/shared"
	spotifyAPI "github.com/zmb3/spotify"
)

// DeletePlaylistTracksUseCase handles the business logic for deleting tracks from a playlist
type DeletePlaylistTracksUseCase struct {
	spotifyRepo shared.SpotifyRepository
	cacheRepo   shared.CacheRepository
}

// NewDeletePlaylistTracksUseCase creates a use case that removes all tracks from a
// specific playlist. It depends on SpotifyRepository for track fetching and deletion,
// and CacheRepository for invalidation.
func NewDeletePlaylistTracksUseCase(
	spotifyRepo shared.SpotifyRepository,
	cacheRepo shared.CacheRepository,
) *DeletePlaylistTracksUseCase {
	return &DeletePlaylistTracksUseCase{
		spotifyRepo: spotifyRepo,
		cacheRepo:   cacheRepo,
	}
}

// Execute removes all tracks from the specified playlist. It fetches the playlist
// tracks (from cache or API), extracts their IDs, deletes them from the playlist,
// and invalidates the playlist tracks cache.
func (uc *DeletePlaylistTracksUseCase) Execute(ctx context.Context, playlistID spotifyAPI.ID) error {
	// Fetch tracks from the playlist, either from cache or Spotify API
	tracks, err := fetchPlaylistTracks(ctx, uc.spotifyRepo, uc.cacheRepo, playlistID)
	if err != nil {
		return err
	}

	if len(tracks) == 0 {
		return nil
	}

	// Extract track IDs and delete them from the playlist
	trackIDs := extractPlaylistTrackIDs(tracks)
	if err := uc.spotifyRepo.DeletePlaylistTracks(ctx, playlistID, trackIDs); err != nil {
		return err
	}

	// Invalidate cache for the playlist and user tracks
	_ = uc.cacheRepo.InvalidatePlaylistTracks(ctx, playlistID)

	return nil
}
