package auth

import (
	"context"

	"github.com/RubenPari/clear-songs/internal/domain/shared"
)

// IsAuthUseCase handles the business logic for checking authentication status
type IsAuthUseCase struct {
	spotifyRepo shared.SpotifyRepository
}

// NewIsAuthUseCase creates a new IsAuthUseCase
func NewIsAuthUseCase(spotifyRepo shared.SpotifyRepository) *IsAuthUseCase {
	return &IsAuthUseCase{
		spotifyRepo: spotifyRepo,
	}
}

// UserInfo represents authenticated user information
type UserInfo struct {
	SpotifyID    string
	DisplayName  string
	Email        string
	ProfileImage string
}

// Execute checks if user is authenticated and returns user info
func (uc *IsAuthUseCase) Execute(ctx context.Context) (*UserInfo, error) {
	// Try to get current user
	user, err := uc.spotifyRepo.GetCurrentUser(ctx)
	if err != nil {
		return nil, err
	}

	profileImage := ""
	if len(user.Images) > 0 {
		profileImage = user.Images[0].URL
	}

	return &UserInfo{
		SpotifyID:    user.ID,
		DisplayName:  user.DisplayName,
		Email:        user.Email,
		ProfileImage: profileImage,
	}, nil
}
