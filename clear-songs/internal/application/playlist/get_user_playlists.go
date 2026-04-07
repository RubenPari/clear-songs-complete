package playlist

import (
	"context"

	"github.com/RubenPari/clear-songs/internal/domain/shared"
	spotifyAPI "github.com/zmb3/spotify"
)

// GetUserPlaylistsUseCase handles the business logic for getting user playlists
type GetUserPlaylistsUseCase struct {
	spotifyRepo shared.SpotifyRepository
	cacheRepo   shared.CacheRepository
}

// NewGetUserPlaylistsUseCase creates a new GetUserPlaylistsUseCase
func NewGetUserPlaylistsUseCase(
	spotifyRepo shared.SpotifyRepository,
	cacheRepo shared.CacheRepository,
) *GetUserPlaylistsUseCase {
	return &GetUserPlaylistsUseCase{
		spotifyRepo: spotifyRepo,
		cacheRepo:   cacheRepo,
	}
}

// Execute retrieves all playlists owned or followed by the user
func (uc *GetUserPlaylistsUseCase) Execute(ctx context.Context) ([]spotifyAPI.SimplePlaylist, error) {
	// Try cache first (if available)
	if uc.cacheRepo != nil {
		// Note: We could cache playlists, but for now we fetch fresh
		// as playlists can change frequently
	}

	// Fetch from API
	playlists, err := uc.spotifyRepo.GetAllUserPlaylists(ctx)
	if err != nil {
		return nil, err
	}

	return playlists, nil
}
