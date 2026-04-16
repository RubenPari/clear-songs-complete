package playlist

import (
	"context"
	"time"

	"github.com/RubenPari/clear-songs/internal/domain/shared"
	spotifyAPI "github.com/zmb3/spotify"
)

const userPlaylistsCacheTTL = 2 * time.Minute

// GetUserPlaylistsUseCase handles the business logic for getting user playlists
type GetUserPlaylistsUseCase struct {
	spotifyRepo shared.SpotifyRepository
	cacheRepo   shared.CacheRepository
}

// Creates get user playlists use case.
func NewGetUserPlaylistsUseCase(
	spotifyRepo shared.SpotifyRepository,
	cacheRepo shared.CacheRepository,
) *GetUserPlaylistsUseCase {
	return &GetUserPlaylistsUseCase{
		spotifyRepo: spotifyRepo,
		cacheRepo:   cacheRepo,
	}
}

// Execute.
func (uc *GetUserPlaylistsUseCase) Execute(ctx context.Context) ([]spotifyAPI.SimplePlaylist, error) {
	cacheKey := "userPlaylists"

	// Check cache first, if available and not expired
	if uc.cacheRepo != nil {
		var cached []spotifyAPI.SimplePlaylist
		found, err := uc.cacheRepo.Get(ctx, cacheKey, &cached)
		if err == nil && found && len(cached) > 0 {
			return cached, nil
		}
	}

	// If not in cache or cache is empty, fetch from Spotify
	playlists, err := uc.spotifyRepo.GetAllUserPlaylists(ctx)
	if err != nil {
		return nil, err
	}

	if uc.cacheRepo != nil {
		_ = uc.cacheRepo.Set(ctx, cacheKey, playlists, userPlaylistsCacheTTL)
	}

	return playlists, nil
}
