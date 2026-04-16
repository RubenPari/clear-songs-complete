package auth

import (
	"golang.org/x/oauth2"
)

// LoginUseCase handles the business logic for initiating OAuth login
type LoginUseCase struct {
	oauthConfig *oauth2.Config
}

// Creates login use case.
func NewLoginUseCase(oauthConfig *oauth2.Config) *LoginUseCase {
	return &LoginUseCase{
		oauthConfig: oauthConfig,
	}
}

// Execute.
func (uc *LoginUseCase) Execute(state string) string {
	if state == "" {
		state = "state"
	}
	return uc.oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
}
