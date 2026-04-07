package track

import (
	"context"
	"time"

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
	cacheRepo   shared.CacheRepository,
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

// getUserTracks retrieves tracks from cache or API
func (uc *DeleteTracksByArtistUseCase) getUserTracks(ctx context.Context) ([]spotifyAPI.SavedTrack, error) {
	// Try cache first (if available)
	if uc.cacheRepo != nil {
		cached, err := uc.cacheRepo.GetUserTracks(ctx)
		if err == nil && cached != nil && len(cached) > 0 {
			return cached, nil
		}
	}

	// Fetch from API
	tracks, err := uc.spotifyRepo.GetAllUserTracks(ctx)
	if err != nil {
		return nil, err
	}

	// Cache for future use (if cache is available)
	if uc.cacheRepo != nil {
		_ = uc.cacheRepo.SetUserTracks(ctx, tracks, 5*time.Minute)
	}

	return tracks, nil
}
