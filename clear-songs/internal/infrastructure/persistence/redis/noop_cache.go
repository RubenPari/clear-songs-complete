package redis

import (
	"context"
	spotifyAPI "github.com/zmb3/spotify"
	"golang.org/x/oauth2"
	"time"
	"github.com/RubenPari/clear-songs/internal/domain/shared"
)

// NoOpCacheRepository is a no-op implementation of CacheRepository used when
// Redis is not available. All operations return nil or empty values.
type NoOpCacheRepository struct{}

// NewNoOpCacheRepository creates a no-op cache repository.
func NewNoOpCacheRepository() shared.CacheRepository {
	return &NoOpCacheRepository{}
}

// SetToken is a no-op.
func (n *NoOpCacheRepository) SetToken(ctx context.Context, token *oauth2.Token) error {
	return nil
}

// GetToken always returns nil.
func (n *NoOpCacheRepository) GetToken(ctx context.Context) (*oauth2.Token, error) {
	return nil, nil
}

// ClearToken is a no-op.
func (n *NoOpCacheRepository) ClearToken(ctx context.Context) error {
	return nil
}

// GetUserTracks always returns nil.
func (n *NoOpCacheRepository) GetUserTracks(ctx context.Context) ([]spotifyAPI.SavedTrack, error) {
	return nil, nil
}

// SetUserTracks is a no-op.
func (n *NoOpCacheRepository) SetUserTracks(ctx context.Context, tracks []spotifyAPI.SavedTrack, ttl time.Duration) error {
	return nil
}

// InvalidateUserTracks is a no-op.
func (n *NoOpCacheRepository) InvalidateUserTracks(ctx context.Context) error {
	return nil
}

// GetPlaylistTracks always returns nil.
func (n *NoOpCacheRepository) GetPlaylistTracks(ctx context.Context, playlistID spotifyAPI.ID) ([]spotifyAPI.PlaylistTrack, error) {
	return nil, nil
}

// SetPlaylistTracks is a no-op.
func (n *NoOpCacheRepository) SetPlaylistTracks(ctx context.Context, playlistID spotifyAPI.ID, tracks []spotifyAPI.PlaylistTrack, ttl time.Duration) error {
	return nil
}

// InvalidatePlaylistTracks is a no-op.
func (n *NoOpCacheRepository) InvalidatePlaylistTracks(ctx context.Context, playlistID spotifyAPI.ID) error {
	return nil
}

// Get always returns (false, nil).
func (n *NoOpCacheRepository) Get(ctx context.Context, key string, target interface{}) (bool, error) {
	return false, nil
}

// Set is a no-op.
func (n *NoOpCacheRepository) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return nil
}

// Delete is a no-op.
func (n *NoOpCacheRepository) Delete(ctx context.Context, key string) error {
	return nil
}

// Ensure NoOpCacheRepository implements CacheRepository interface
var _ shared.CacheRepository = (*NoOpCacheRepository)(nil)
