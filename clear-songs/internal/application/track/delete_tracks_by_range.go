package track

import (
	"context"

	"github.com/RubenPari/clear-songs/internal/domain/shared"
	spotifyAPI "github.com/zmb3/spotify"
)

// DeleteTracksByRangeUseCase handles the business logic for deleting tracks by range.
type DeleteTracksByRangeUseCase struct {
	spotifyRepo shared.SpotifyRepository
	cacheRepo   shared.CacheRepository
}

// NewDeleteTracksByRangeUseCase creates a new DeleteTracksByRangeUseCase.
func NewDeleteTracksByRangeUseCase(
	spotifyRepo shared.SpotifyRepository,
	cacheRepo shared.CacheRepository,
) *DeleteTracksByRangeUseCase {
	return &DeleteTracksByRangeUseCase{
		spotifyRepo: spotifyRepo,
		cacheRepo:   cacheRepo,
	}
}

// Execute deletes all saved tracks whose primary artist falls within the track-count range [min, max].
func (uc *DeleteTracksByRangeUseCase) Execute(ctx context.Context, min, max int) error {
	tracks, err := getUserTracks(ctx, uc.spotifyRepo, uc.cacheRepo)
	if err != nil {
		return err
	}

	artistMap := groupTracksByPrimaryArtist(tracks)

	var trackIDsToDelete []spotifyAPI.ID
	for _, artist := range artistMap {
		if !passesRangeFilter(artist.count, min, max) {
			continue
		}

		artistTrackIDs, err := uc.spotifyRepo.GetTrackIDsByArtist(ctx, spotifyAPI.ID(artist.id), tracks)
		if err != nil {
			return err
		}

		trackIDsToDelete = append(trackIDsToDelete, artistTrackIDs...)
	}

	if len(trackIDsToDelete) == 0 {
		return nil
	}

	if err := uc.spotifyRepo.DeleteTracksFromLibrary(ctx, trackIDsToDelete); err != nil {
		return err
	}

	_ = uc.cacheRepo.InvalidateUserTracks(ctx)
	return nil
}
