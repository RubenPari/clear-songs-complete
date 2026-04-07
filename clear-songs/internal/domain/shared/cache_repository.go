package shared

import (
	"context"
	"time"
	spotifyAPI "github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

// CacheRepository defines the interface for caching operations
type CacheRepository interface {
	// Token operations
	SetToken(ctx context.Context, token *oauth2.Token) error
	GetToken(ctx context.Context) (*oauth2.Token, error)
	ClearToken(ctx context.Context) error
	
	// User tracks cache
	GetUserTracks(ctx context.Context) ([]spotifyAPI.SavedTrack, error)
	SetUserTracks(ctx context.Context, tracks []spotifyAPI.SavedTrack, ttl time.Duration) error
	InvalidateUserTracks(ctx context.Context) error
	
	// Playlist tracks cache
	GetPlaylistTracks(ctx context.Context, playlistID spotifyAPI.ID) ([]spotifyAPI.PlaylistTrack, error)
	SetPlaylistTracks(ctx context.Context, playlistID spotifyAPI.ID, tracks []spotifyAPI.PlaylistTrack, ttl time.Duration) error
	InvalidatePlaylistTracks(ctx context.Context, playlistID spotifyAPI.ID) error
	
	// Generic cache operations
	Get(ctx context.Context, key string, target interface{}) (bool, error)
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
}
