package auth

import (
	"context"

	"github.com/RubenPari/clear-songs/internal/domain/shared"
)

// LogoutUseCase handles the business logic for user logout.
type LogoutUseCase struct {
	spotifyRepo shared.SpotifyRepository
	cacheRepo   shared.CacheRepository
}

// NewLogoutUseCase creates a LogoutUseCase.
func NewLogoutUseCase(
	spotifyRepo shared.SpotifyRepository,
	cacheRepo shared.CacheRepository,
) *LogoutUseCase {
	return &LogoutUseCase{
		spotifyRepo: spotifyRepo,
		cacheRepo:   cacheRepo,
	}
}

// Execute clears the current OAuth token from the repository and cache.
func (uc *LogoutUseCase) Execute(ctx context.Context) error {
	_ = uc.spotifyRepo.SetAccessToken(nil)

	if uc.cacheRepo != nil {
		_ = uc.cacheRepo.ClearToken(ctx)
	}

	return nil
}
