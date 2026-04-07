package auth

import (
	"golang.org/x/oauth2"
)

// LoginUseCase handles the business logic for initiating OAuth login
type LoginUseCase struct {
	oauthConfig *oauth2.Config
}

// NewLoginUseCase creates a new LoginUseCase
func NewLoginUseCase(oauthConfig *oauth2.Config) *LoginUseCase {
	return &LoginUseCase{
		oauthConfig: oauthConfig,
	}
}

// Execute generates the OAuth authorization URL
func (uc *LoginUseCase) Execute(state string) string {
	if state == "" {
		state = "state"
	}
	return uc.oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
}
