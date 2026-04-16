package track

import (
	"context"
	"fmt"

	"github.com/RubenPari/clear-songs/internal/domain/shared"
	spotifyAPI "github.com/zmb3/spotify"
	"go.uber.org/zap"
)

// DeleteTrackUseCase handles the business logic for deleting a single track
type DeleteTrackUseCase struct {
	spotifyRepo  shared.SpotifyRepository
	cacheRepo    shared.CacheRepository
	databaseRepo shared.DatabaseRepository
}

// Creates delete track use case.
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

// Execute.
func (uc *DeleteTrackUseCase) Execute(ctx context.Context, trackID spotifyAPI.ID) error {
	// 1. Get track details for backup
	track, err := uc.spotifyRepo.GetTrack(ctx, trackID)
	if err != nil {
		return fmt.Errorf("get track details: %w", err)
	}

	// 2. Save backup to database
	if uc.databaseRepo != nil {
		if err := uc.databaseRepo.SaveFullTracksBackup([]spotifyAPI.FullTrack{*track}); err != nil {
			zap.L().Warn("failed to backup track before deletion", zap.String("track_id", trackID.String()), zap.Error(err))
		}
	}

	// 3. Delete track from library
	if err := uc.spotifyRepo.DeleteTracksFromLibrary(ctx, []spotifyAPI.ID{trackID}); err != nil {
		return fmt.Errorf("delete track from library: %w", err)
	}

	// 4. Invalidate cache
	if uc.cacheRepo != nil {
		_ = uc.cacheRepo.InvalidateUserTracks(ctx)
	}

	return nil
}
