// Package constants holds shared domain-wide constants and configuration
// values used across the application and infrastructure layers.
package constants

// Scopes defines the full set of Spotify OAuth 2.0 scopes requested during
// the authorisation flow. These scopes grant access to playlist management,
// library read/write, user profile, playback state, and followed-artists
// endpoints required by the application.
var Scopes = []string{"playlist-read-private", "playlist-read-collaborative", "playlist-modify-public", "playlist-modify-private", "user-library-read", "user-library-modify", "user-read-private", "user-read-email", "user-read-playback-state", "user-modify-playback-state", "user-read-currently-playing", "user-read-recently-played", "user-top-read", "user-follow-read", "user-follow-modify"}
