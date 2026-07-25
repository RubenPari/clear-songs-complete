package playlist

import (
	"context"
	"time"

	"github.com/RubenPari/clear-songs/internal/domain/shared"
	spotifyAPI "github.com/zmb3/spotify"
)

// userPlaylistsCacheTTL is the duration for which the user's playlists are cached
// in Redis before being re-fetched from the Spotify API.
const userPlaylistsCacheTTL = 2 * time.Minute

// GetUserPlaylistsUseCase handles the business logic for getting user playlists
type GetUserPlaylistsUseCase struct {
	spotifyRepo shared.SpotifyRepository
	cacheRepo   shared.CacheRepository
}

// NewGetUserPlaylistsUseCase creates a use case that retrieves all playlists owned
// by the authenticated user. It depends on SpotifyRepository for playlist fetching
// and CacheRepository for caching.
func NewGetUserPlaylistsUseCase(
	spotifyRepo shared.SpotifyRepository,
	cacheRepo shared.CacheRepository,
) *GetUserPlaylistsUseCase {
	return &GetUserPlaylistsUseCase{
		spotifyRepo: spotifyRepo,
		cacheRepo:   cacheRepo,
	}
}

// Execute retrieves all playlists owned by the authenticated user. It checks the
// cache first, and on a miss it fetches all playlists from Spotify and caches the
// result with a 2-minute TTL.
func (uc *GetUserPlaylistsUseCase) Execute(ctx context.Context) ([]spotifyAPI.SimplePlaylist, error) {
	cacheKey := shared.CacheKeyUserPlaylists

	// Check cache first, if available and not expired
	var cached []spotifyAPI.SimplePlaylist
	found, err := uc.cacheRepo.Get(ctx, cacheKey, &cached)
	if err == nil && found && len(cached) > 0 {
		return cached, nil
	}

	// If not in cache or cache is empty, fetch from Spotify
	playlists, err := uc.spotifyRepo.GetAllUserPlaylists(ctx)
	if err != nil {
		return nil, err
	}

	// Store the fetched playlists in cache for future requests
	_ = uc.cacheRepo.Set(ctx, cacheKey, playlists, userPlaylistsCacheTTL)

	return playlists, nil
}
