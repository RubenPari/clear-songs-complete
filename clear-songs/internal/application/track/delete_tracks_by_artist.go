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

// NewDeleteTracksByArtistUseCase creates a use case that removes all saved tracks
// by a specific artist from the user's library. It depends on SpotifyRepository
// for track fetching and deletion, and CacheRepository for invalidation.
func NewDeleteTracksByArtistUseCase(
	spotifyRepo shared.SpotifyRepository,
	cacheRepo shared.CacheRepository,
) *DeleteTracksByArtistUseCase {
	return &DeleteTracksByArtistUseCase{
		spotifyRepo: spotifyRepo,
		cacheRepo:   cacheRepo,
	}
}

// Execute removes all saved tracks whose primary artist matches the given artistID.
// It fetches the user's tracks (from cache or API), filters by artist, deletes the
// matching tracks from the library, and invalidates the user tracks cache.
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
	_ = uc.cacheRepo.InvalidateUserTracks(ctx)

	return nil
}

// getUserTracks delegates to the shared helper that reads from cache or fetches
// from the Spotify API.
func (uc *DeleteTracksByArtistUseCase) getUserTracks(ctx context.Context) ([]spotifyAPI.SavedTrack, error) {
	return getUserTracks(ctx, uc.spotifyRepo, uc.cacheRepo)
}
