package track

import (
	"context"

	"github.com/RubenPari/clear-songs/internal/domain/shared"
	spotifyAPI "github.com/zmb3/spotify"
)

// DeleteTrackUseCase handles the business logic for deleting a single track
type DeleteTrackUseCase struct {
	spotifyRepo  shared.SpotifyRepository
	cacheRepo    shared.CacheRepository
	databaseRepo shared.DatabaseRepository
}

// NewDeleteTrackUseCase creates a new DeleteTrackUseCase
func NewDeleteTrackUseCase(
	spotifyRepo shared.SpotifyRepository,
	cacheRepo shared.CacheRepository,
	databaseRepo shared.DatabaseRepository,
) *DeleteTrackUseCase {
	return &DeleteTrackUseCase{
		spotifyRepo:  spotifyRepo,
		cacheRepo:    cacheRepo,
		databaseRepo: databaseRepo,
	}
}

// Execute deletes a single track from the user's library
func (uc *DeleteTrackUseCase) Execute(ctx context.Context, trackID spotifyAPI.ID) error {
	// 1. Get track details for backup
	track, err := uc.spotifyRepo.GetTrack(ctx, trackID)
	if err != nil {
		return err
	}

	// 2. Save backup to database
	if err := uc.databaseRepo.SaveFullTracksBackup([]spotifyAPI.FullTrack{*track}); err != nil {
		// Log error but continue
	}

	// 3. Delete track from library
	if err := uc.spotifyRepo.DeleteTracksFromLibrary(ctx, []spotifyAPI.ID{trackID}); err != nil {
		return err
	}

	// 4. Invalidate cache
	if uc.cacheRepo != nil {
		_ = uc.cacheRepo.InvalidateUserTracks(ctx)
	}

	return nil
}
