// Package auth defines the domain entity for authenticated users and provides
// conversion from Spotify API user profiles.
package auth

import spotifyAPI "github.com/zmb3/spotify"

// User represents the currently authenticated Spotify user.
type User struct {
	ID           string
	DisplayName  string
	Email        string
	ProfileImage string
}

// NewUserFromSpotify builds a domain User from a Spotify private user profile.
// The images slice is expected to be ordered by preference (e.g., largest first).
func NewUserFromSpotify(user *spotifyAPI.PrivateUser) *User {
	if user == nil {
		return nil
	}

	profileImage := ""
	if len(user.Images) > 0 {
		profileImage = user.Images[0].URL
	}

	return &User{
		ID:           user.ID,
		DisplayName:  user.DisplayName,
		Email:        user.Email,
		ProfileImage: profileImage,
	}
}
