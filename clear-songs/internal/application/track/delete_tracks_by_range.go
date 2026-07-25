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

// NewDeleteTracksByRangeUseCase creates a use case that removes all saved tracks
// whose primary artist has a track count within the specified range. It depends on
// SpotifyRepository for track fetching and deletion, and CacheRepository for invalidation.
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
	// Fetch all saved tracks from the user's library, either from cache or Spotify API
	tracks, err := getUserTracks(ctx, uc.spotifyRepo, uc.cacheRepo)
	if err != nil {
		return err
	}

	// Group tracks by their primary artist and filter based on the specified range
	artistMap := groupTracksByPrimaryArtist(tracks)

	// Collect track IDs to delete based on the range filter
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

	// Delete the filtered tracks from the user's library
	if err := uc.spotifyRepo.DeleteTracksFromLibrary(ctx, trackIDsToDelete); err != nil {
		return err
	}

	// Invalidate cache for the user's saved tracks
	_ = uc.cacheRepo.InvalidateUserTracks(ctx)

	return nil
}
