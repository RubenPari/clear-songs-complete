package auth

import (
	"context"
	"fmt"
	"os"

	"github.com/RubenPari/clear-songs/internal/domain/shared"
	"golang.org/x/oauth2"
)

const defaultFrontendURL = "http://127.0.0.1:4200"

// CallbackUseCase handles the business logic for OAuth callback.
type CallbackUseCase struct {
	oauthConfig *oauth2.Config
	spotifyRepo shared.SpotifyRepository
	cacheRepo   shared.CacheRepository
}

// NewCallbackUseCase creates a CallbackUseCase.
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

// Execute exchanges the authorization code for a token, persists it and returns the frontend redirect URL.
func (uc *CallbackUseCase) Execute(ctx context.Context, code string) (string, error) {
	token, err := uc.oauthConfig.Exchange(ctx, code)
	if err != nil {
		return "", fmt.Errorf("oauth2 exchange: %w", err)
	}

	if uc.cacheRepo != nil {
		if err := uc.cacheRepo.SetToken(ctx, token); err != nil {
			return "", fmt.Errorf("cache set token: %w", err)
		}
	}

	if err := uc.spotifyRepo.SetAccessToken(token); err != nil {
		return "", fmt.Errorf("spotify set access token: %w", err)
	}

	if _, err := uc.spotifyRepo.GetCurrentUser(ctx); err != nil {
		return "", fmt.Errorf("spotify get current user: %w", err)
	}

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = defaultFrontendURL
	}

	return frontendURL + "/callback", nil
}
