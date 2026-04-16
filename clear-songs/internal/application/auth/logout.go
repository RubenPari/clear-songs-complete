package auth

import (
	"context"

	"github.com/RubenPari/clear-songs/internal/domain/shared"
)

// LogoutUseCase handles the business logic for user logout
type LogoutUseCase struct {
	spotifyRepo shared.SpotifyRepository
	cacheRepo   shared.CacheRepository
}

// Creates logout use case.
func NewLogoutUseCase(
	spotifyRepo shared.SpotifyRepository,
	cacheRepo shared.CacheRepository,
) *LogoutUseCase {
	return &LogoutUseCase{
		spotifyRepo: spotifyRepo,
		cacheRepo:   cacheRepo,
	}
}

// Execute.
func (uc *LogoutUseCase) Execute(ctx context.Context) error {
	// Clear token from Spotify repository
	_ = uc.spotifyRepo.SetAccessToken(nil)

	// Clear token from cache
	if uc.cacheRepo != nil {
		_ = uc.cacheRepo.ClearToken(ctx)
	}

	return nil
}
