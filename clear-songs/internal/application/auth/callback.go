package auth

import (
	"context"
	"fmt"
	"os"

	"github.com/RubenPari/clear-songs/internal/domain/shared"
	"golang.org/x/oauth2"
)

const defaultFrontendURL = "http://127.0.0.1:4200"

// CallbackUseCase handles the business logic for OAuth callback
type CallbackUseCase struct {
	oauthConfig *oauth2.Config
	spotifyRepo shared.SpotifyRepository
	cacheRepo   shared.CacheRepository
}

// Creates callback use case.
func NewCallbackUseCase(
	oauthConfig *oauth2.Config,
	spotifyRepo shared.SpotifyRepository,
	cacheRepo shared.CacheRepository,
) *CallbackUseCase {
	return &CallbackUseCase{
		oauthConfig: oauthConfig,
		spotifyRepo: spotifyRepo,
		cacheRepo:   cacheRepo,
	}
}

// Execute.
func (uc *CallbackUseCase) Execute(ctx context.Context, code string) (string, error) {
	// 1. Exchange code for token
	token, err := uc.oauthConfig.Exchange(ctx, code)
	if err != nil {
		return "", fmt.Errorf("oauth2 exchange: %w", err)
	}

	// 2. Save token to cache
	if uc.cacheRepo != nil {
		if err := uc.cacheRepo.SetToken(ctx, token); err != nil {
			return "", fmt.Errorf("cache set token: %w", err)
		}
	}

	// 3. Set token in Spotify repository
	if err := uc.spotifyRepo.SetAccessToken(token); err != nil {
		return "", fmt.Errorf("spotify set access token: %w", err)
	}

	// 4. Verify authentication by getting current user
	_, err = uc.spotifyRepo.GetCurrentUser(ctx)
	if err != nil {
		return "", fmt.Errorf("spotify get current user: %w", err)
	}

	// 5. Get frontend URL
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = defaultFrontendURL
	}

	return frontendURL + "/callback", nil
}
