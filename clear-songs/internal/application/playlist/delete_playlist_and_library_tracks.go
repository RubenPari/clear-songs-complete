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

// Creates delete playlist and library tracks use case.
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

// Execute.
func (uc *DeletePlaylistAndLibraryTracksUseCase) Execute(ctx context.Context, playlistID spotifyAPI.ID) error {
	tracks, err := fetchPlaylistTracks(ctx, uc.spotifyRepo, uc.cacheRepo, playlistID)
	if err != nil {
		return err
	}

	if len(tracks) == 0 {
		return nil
	}

	if err := uc.databaseRepo.SaveTracksBackup(tracks); err != nil {
		zap.L().Warn("failed to save playlist tracks backup", zap.Error(err))
	}

	trackIDs := extractPlaylistTrackIDs(tracks)

	if err := uc.spotifyRepo.DeletePlaylistTracks(ctx, playlistID, trackIDs); err != nil {
		return err
	}

	if err := uc.spotifyRepo.DeleteTracksFromLibrary(ctx, trackIDs); err != nil {
		return err
	}

	_ = uc.cacheRepo.InvalidatePlaylistTracks(ctx, playlistID)
	_ = uc.cacheRepo.InvalidateUserTracks(ctx)

	return nil
}
