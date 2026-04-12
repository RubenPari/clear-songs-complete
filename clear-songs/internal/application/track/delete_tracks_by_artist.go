package track

import (
	"context"

	"github.com/RubenPari/clear-songs/internal/domain/shared"
	spotifyAPI "github.com/zmb3/spotify"
)

// DeleteTracksByArtistUseCase handles the business logic for deleting tracks by artist
type DeleteTracksByArtistUseCase struct {
	spotifyRepo shared.SpotifyRepository
	cacheRepo   shared.CacheRepository
}

// NewDeleteTracksByArtistUseCase creates a new DeleteTracksByArtistUseCase
func NewDeleteTracksByArtistUseCase(
	spotifyRepo shared.SpotifyRepository,
	cacheRepo shared.CacheRepository,
) *DeleteTracksByArtistUseCase {
	return &DeleteTracksByArtistUseCase{
		spotifyRepo: spotifyRepo,
		cacheRepo:   cacheRepo,
	}
}

// Execute deletes all tracks from a specific artist
func (uc *DeleteTracksByArtistUseCase) Execute(ctx context.Context, artistID spotifyAPI.ID) error {
	// 1. Get user tracks (from cache or API)
	tracks, err := uc.getUserTracks(ctx)
	if err != nil {
		return err
	}

	// 2. Filter tracks by artist
	trackIDs, err := uc.spotifyRepo.GetTrackIDsByArtist(ctx, artistID, tracks)
	if err != nil {
		return err
	}

	if len(trackIDs) == 0 {
		return nil // No tracks to delete
	}

	// 3. Delete tracks from library
	if err := uc.spotifyRepo.DeleteTracksFromLibrary(ctx, trackIDs); err != nil {
		return err
	}

	// 4. Invalidate cache
	if uc.cacheRepo != nil {
		_ = uc.cacheRepo.InvalidateUserTracks(ctx)
	}

	return nil
}

// getUserTracks retrieves tracks from cache or API.
func (uc *DeleteTracksByArtistUseCase) getUserTracks(ctx context.Context) ([]spotifyAPI.SavedTrack, error) {
	return getUserTracks(ctx, uc.spotifyRepo, uc.cacheRepo)
}
