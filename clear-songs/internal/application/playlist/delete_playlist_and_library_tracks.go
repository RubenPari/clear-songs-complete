package playlist

import (
	"context"
	"time"

	"github.com/RubenPari/clear-songs/internal/domain/shared"
	spotifyAPI "github.com/zmb3/spotify"
)

// DeletePlaylistAndLibraryTracksUseCase handles the business logic for deleting tracks from both playlist and library
type DeletePlaylistAndLibraryTracksUseCase struct {
	spotifyRepo    shared.SpotifyRepository
	cacheRepo      shared.CacheRepository
	databaseRepo   shared.DatabaseRepository
	deletePlaylistUC *DeletePlaylistTracksUseCase
}

// NewDeletePlaylistAndLibraryTracksUseCase creates a new DeletePlaylistAndLibraryTracksUseCase
func NewDeletePlaylistAndLibraryTracksUseCase(
	spotifyRepo shared.SpotifyRepository,
	cacheRepo shared.CacheRepository,
	databaseRepo shared.DatabaseRepository,
	deletePlaylistUC *DeletePlaylistTracksUseCase,
) *DeletePlaylistAndLibraryTracksUseCase {
	return &DeletePlaylistAndLibraryTracksUseCase{
		spotifyRepo:      spotifyRepo,
		cacheRepo:        cacheRepo,
		databaseRepo:     databaseRepo,
		deletePlaylistUC: deletePlaylistUC,
	}
}

// Execute deletes tracks from both playlist and user library
func (uc *DeletePlaylistAndLibraryTracksUseCase) Execute(ctx context.Context, playlistID spotifyAPI.ID) error {
	// 1. Get playlist tracks (from cache or API)
	tracks, err := uc.getPlaylistTracks(ctx, playlistID)
	if err != nil {
		return err
	}

	if len(tracks) == 0 {
		return nil // No tracks to delete
	}

	// 2. Save backup to database
	if err := uc.databaseRepo.SaveTracksBackup(tracks); err != nil {
		// Log error but continue (backup is optional)
	}

	// 3. Delete tracks from playlist (reuse existing use case)
	if err := uc.deletePlaylistUC.Execute(ctx, playlistID); err != nil {
		return err
	}

	// 4. Convert tracks to IDs for library deletion
	trackIDs := make([]spotifyAPI.ID, 0, len(tracks))
	for _, track := range tracks {
		trackIDs = append(trackIDs, track.Track.ID)
	}

	// 5. Delete tracks from user library
	if err := uc.spotifyRepo.DeleteTracksFromLibrary(ctx, trackIDs); err != nil {
		return err
	}

	// 6. Invalidate cache (both playlist and user data)
	if uc.cacheRepo != nil {
		_ = uc.cacheRepo.InvalidatePlaylistTracks(ctx, playlistID)
		_ = uc.cacheRepo.InvalidateUserTracks(ctx)
	}

	return nil
}

// getPlaylistTracks retrieves tracks from cache or API
func (uc *DeletePlaylistAndLibraryTracksUseCase) getPlaylistTracks(ctx context.Context, playlistID spotifyAPI.ID) ([]spotifyAPI.PlaylistTrack, error) {
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
