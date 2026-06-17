package auth

import (
	"context"

	"github.com/RubenPari/clear-songs/internal/domain/auth"
	"github.com/RubenPari/clear-songs/internal/domain/shared"
)

// IsAuthUseCase handles the business logic for checking authentication status.
type IsAuthUseCase struct {
	spotifyRepo shared.SpotifyRepository
}

// NewIsAuthUseCase creates an IsAuthUseCase.
func NewIsAuthUseCase(spotifyRepo shared.SpotifyRepository) *IsAuthUseCase {
	return &IsAuthUseCase{
		spotifyRepo: spotifyRepo,
	}
}

// Execute returns the authenticated user, or an error if the user is not authenticated.
func (uc *IsAuthUseCase) Execute(ctx context.Context) (*auth.User, error) {
	user, err := uc.spotifyRepo.GetCurrentUser(ctx)
	if err != nil {
		return nil, err
	}

	return auth.NewUserFromSpotify(user), nil
}
