package playlist

import (
	"context"
	"time"

	"github.com/RubenPari/clear-songs/internal/domain/shared"
	spotifyAPI "github.com/zmb3/spotify"
)

// DeletePlaylistTracksUseCase handles the business logic for deleting tracks from a playlist
type DeletePlaylistTracksUseCase struct {
	spotifyRepo shared.SpotifyRepository
	cacheRepo   shared.CacheRepository
}

// NewDeletePlaylistTracksUseCase creates a new DeletePlaylistTracksUseCase
func NewDeletePlaylistTracksUseCase(
	spotifyRepo shared.SpotifyRepository,
	cacheRepo shared.CacheRepository,
) *DeletePlaylistTracksUseCase {
	return &DeletePlaylistTracksUseCase{
		spotifyRepo: spotifyRepo,
		cacheRepo:   cacheRepo,
	}
}

// Execute deletes all tracks from a playlist
func (uc *DeletePlaylistTracksUseCase) Execute(ctx context.Context, playlistID spotifyAPI.ID) error {
	// 1. Get playlist tracks (from cache or API)
	tracks, err := uc.getPlaylistTracks(ctx, playlistID)
	if err != nil {
		return err
	}

	if len(tracks) == 0 {
		return nil // No tracks to delete
	}

	// 2. Convert tracks to IDs
	trackIDs := make([]spotifyAPI.ID, 0, len(tracks))
	for _, track := range tracks {
		trackIDs = append(trackIDs, track.Track.ID)
	}

	// 3. Delete tracks from playlist
	if err := uc.spotifyRepo.DeletePlaylistTracks(ctx, playlistID, trackIDs); err != nil {
		return err
	}

	// 4. Invalidate cache
	if uc.cacheRepo != nil {
		_ = uc.cacheRepo.InvalidatePlaylistTracks(ctx, playlistID)
	}

	return nil
}

// getPlaylistTracks retrieves tracks from cache or API
func (uc *DeletePlaylistTracksUseCase) getPlaylistTracks(ctx context.Context, playlistID spotifyAPI.ID) ([]spotifyAPI.PlaylistTrack, error) {
	// Try cache first (if available)
	if uc.cacheRepo != nil {
		cached, err := uc.cacheRepo.GetPlaylistTracks(ctx, playlistID)
		if err == nil && cached != nil && len(cached) > 0 {
			return cached, nil
		}
	}

	// Fetch from API
	tracks, err := uc.spotifyRepo.GetAllPlaylistTracks(ctx, playlistID)
	if err != nil {
		return nil, err
	}

	// Cache for future use (if cache is available)
	if uc.cacheRepo != nil {
		_ = uc.cacheRepo.SetPlaylistTracks(ctx, playlistID, tracks, 5*time.Minute)
	}

	return tracks, nil
}
