package playlist

import (
	"context"

	"github.com/RubenPari/clear-songs/internal/domain/shared"
	spotifyAPI "github.com/zmb3/spotify"
	"go.uber.org/zap"
)

// DeletePlaylistAndLibraryTracksUseCase handles the business logic for deleting tracks from both playlist and library
type DeletePlaylistAndLibraryTracksUseCase struct {
	spotifyRepo  shared.SpotifyRepository
	cacheRepo    shared.CacheRepository
	databaseRepo shared.DatabaseRepository
}

// NewDeletePlaylistAndLibraryTracksUseCase creates a use case that removes all tracks
// from a playlist and also removes them from the user's saved library. It backs up
// tracks to Postgres before deletion. It depends on SpotifyRepository for track operations,
// CacheRepository for invalidation, and DatabaseRepository for backup storage.
func NewDeletePlaylistAndLibraryTracksUseCase(
	spotifyRepo shared.SpotifyRepository,
	cacheRepo shared.CacheRepository,
	databaseRepo shared.DatabaseRepository,
) *DeletePlaylistAndLibraryTracksUseCase {
	return &DeletePlaylistAndLibraryTracksUseCase{
		spotifyRepo:  spotifyRepo,
		cacheRepo:    cacheRepo,
		databaseRepo: databaseRepo,
	}
}

// Execute removes all tracks from the specified playlist and also removes them from
// the user's saved library. It fetches the playlist tracks, saves a backup to the
// database (logging a warning on failure), deletes from both the playlist and library,
// and invalidates both caches.
func (uc *DeletePlaylistAndLibraryTracksUseCase) Execute(ctx context.Context, playlistID spotifyAPI.ID) error {
	// Fetch tracks from the playlist, either from cache or Spotify API
	tracks, err := fetchPlaylistTracks(ctx, uc.spotifyRepo, uc.cacheRepo, playlistID)
	if err != nil {
		return err
	}

	if len(tracks) == 0 {
		return nil
	}

	// Save a backup of the tracks to the database before deletion
	if err := uc.databaseRepo.SaveTracksBackup(tracks); err != nil {
		zap.L().Warn("failed to save playlist tracks backup", zap.Error(err))
	}

	trackIDs := extractPlaylistTrackIDs(tracks)

	// Delete tracks from the playlist and the user's library
	if err := uc.spotifyRepo.DeletePlaylistTracks(ctx, playlistID, trackIDs); err != nil {
		return err
	}

	// Delete tracks from the user's library
	if err := uc.spotifyRepo.DeleteTracksFromLibrary(ctx, trackIDs); err != nil {
		return err
	}

	// Invalidate cache for the playlist and user tracks
	_ = uc.cacheRepo.InvalidatePlaylistTracks(ctx, playlistID)
	_ = uc.cacheRepo.InvalidateUserTracks(ctx)

	return nil
}
