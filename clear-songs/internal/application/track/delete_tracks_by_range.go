package track

import (
	"context"

	"github.com/RubenPari/clear-songs/internal/domain/shared"
	spotifyAPI "github.com/zmb3/spotify"
)

// DeleteTracksByRangeUseCase handles the business logic for deleting tracks by range
type DeleteTracksByRangeUseCase struct {
	spotifyRepo       shared.SpotifyRepository
	cacheRepo         shared.CacheRepository
	getTrackSummaryUC *GetTrackSummaryUseCase
	deleteByArtistUC  *DeleteTracksByArtistUseCase
}

// NewDeleteTracksByRangeUseCase creates a new DeleteTracksByRangeUseCase
func NewDeleteTracksByRangeUseCase(
	spotifyRepo shared.SpotifyRepository,
	cacheRepo shared.CacheRepository,
	getTrackSummaryUC *GetTrackSummaryUseCase,
	deleteByArtistUC *DeleteTracksByArtistUseCase,
) *DeleteTracksByRangeUseCase {
	return &DeleteTracksByRangeUseCase{
		spotifyRepo:       spotifyRepo,
		cacheRepo:         cacheRepo,
		getTrackSummaryUC: getTrackSummaryUC,
		deleteByArtistUC:  deleteByArtistUC,
	}
}

// Execute deletes tracks within a count range
func (uc *DeleteTracksByRangeUseCase) Execute(ctx context.Context, min, max int) error {
	// 1. Get track summary filtered by range (no genre filter for deletion)
	summary, err := uc.getTrackSummaryUC.Execute(ctx, min, max, "")
	if err != nil {
		return err
	}

	// 2. Delete tracks for each artist in the summary
	deletionsOccurred := false
	for _, artist := range summary {
		if err := uc.deleteByArtistUC.Execute(ctx, spotifyAPI.ID(artist.ID)); err != nil {
			return err
		}
		deletionsOccurred = true
	}

	// 3. Invalidate cache if deletions occurred
	if deletionsOccurred && uc.cacheRepo != nil {
		_ = uc.cacheRepo.InvalidateUserTracks(ctx)
	}

	return nil
}
