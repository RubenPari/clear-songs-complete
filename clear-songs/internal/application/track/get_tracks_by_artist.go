package track

import (
	"context"

	"github.com/RubenPari/clear-songs/internal/domain/shared"
	spotifyAPI "github.com/zmb3/spotify"
)

// GetTracksByArtistUseCase handles the business logic for getting tracks by artist
type GetTracksByArtistUseCase struct {
	spotifyRepo shared.SpotifyRepository
	cacheRepo   shared.CacheRepository
}

// NewGetTracksByArtistUseCase creates a new GetTracksByArtistUseCase
func NewGetTracksByArtistUseCase(
	spotifyRepo shared.SpotifyRepository,
	cacheRepo shared.CacheRepository,
) *GetTracksByArtistUseCase {
	return &GetTracksByArtistUseCase{
		spotifyRepo: spotifyRepo,
		cacheRepo:   cacheRepo,
	}
}

// Execute retrieves all tracks from a specific artist
func (uc *GetTracksByArtistUseCase) Execute(ctx context.Context, artistID spotifyAPI.ID) ([]spotifyAPI.SavedTrack, error) {
	// 1. Get user tracks (from cache or API)
	tracks, err := uc.getUserTracks(ctx)
	if err != nil {
		return nil, err
	}

	// 2. Filter tracks by artist
	filteredTracks, err := uc.spotifyRepo.GetTracksByArtist(ctx, artistID, tracks)
	if err != nil {
		return nil, err
	}

	return filteredTracks, nil
}

// getUserTracks retrieves tracks from cache or API.
func (uc *GetTracksByArtistUseCase) getUserTracks(ctx context.Context) ([]spotifyAPI.SavedTrack, error) {
	return getUserTracks(ctx, uc.spotifyRepo, uc.cacheRepo)
}
