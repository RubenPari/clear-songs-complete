package redis

import (
	"context"
	spotifyAPI "github.com/zmb3/spotify"
	"golang.org/x/oauth2"
	"time"
	"github.com/RubenPari/clear-songs/internal/domain/shared"
)

// NoOpCacheRepository is a no-op implementation of CacheRepository
// Used when Redis is not available
type NoOpCacheRepository struct{}

// NewNoOpCacheRepository creates a new no-op cache repository
func NewNoOpCacheRepository() shared.CacheRepository {
	return &NoOpCacheRepository{}
}

func (n *NoOpCacheRepository) SetToken(ctx context.Context, token *oauth2.Token) error {
	return nil // No-op
}

func (n *NoOpCacheRepository) GetToken(ctx context.Context) (*oauth2.Token, error) {
	return nil, nil // No token found
}

func (n *NoOpCacheRepository) ClearToken(ctx context.Context) error {
	return nil // No-op
}

func (n *NoOpCacheRepository) GetUserTracks(ctx context.Context) ([]spotifyAPI.SavedTrack, error) {
	return nil, nil // No cache
}

func (n *NoOpCacheRepository) SetUserTracks(ctx context.Context, tracks []spotifyAPI.SavedTrack, ttl time.Duration) error {
	return nil // No-op
}

func (n *NoOpCacheRepository) InvalidateUserTracks(ctx context.Context) error {
	return nil // No-op
}

func (n *NoOpCacheRepository) GetPlaylistTracks(ctx context.Context, playlistID spotifyAPI.ID) ([]spotifyAPI.PlaylistTrack, error) {
	return nil, nil // No cache
}

func (n *NoOpCacheRepository) SetPlaylistTracks(ctx context.Context, playlistID spotifyAPI.ID, tracks []spotifyAPI.PlaylistTrack, ttl time.Duration) error {
	return nil // No-op
}

func (n *NoOpCacheRepository) InvalidatePlaylistTracks(ctx context.Context, playlistID spotifyAPI.ID) error {
	return nil // No-op
}

func (n *NoOpCacheRepository) Get(ctx context.Context, key string, target interface{}) (bool, error) {
	return false, nil // Not found
}

func (n *NoOpCacheRepository) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return nil // No-op
}

func (n *NoOpCacheRepository) Delete(ctx context.Context, key string) error {
	return nil // No-op
}

// Ensure NoOpCacheRepository implements CacheRepository interface
var _ shared.CacheRepository = (*NoOpCacheRepository)(nil)
