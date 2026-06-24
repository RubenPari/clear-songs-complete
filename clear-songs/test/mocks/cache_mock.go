package mocks

import (
	"context"
	"time"

	"github.com/RubenPari/clear-songs/internal/domain/shared"
	"github.com/stretchr/testify/mock"
	spotifyAPI "github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

// MockCacheRepository is a mock implementation of CacheRepository
type MockCacheRepository struct {
	mock.Mock
}

func (m *MockCacheRepository) Get(ctx context.Context, key string, target interface{}) (bool, error) {
	args := m.Called(ctx, key, target)
	return args.Bool(0), args.Error(1)
}

func (m *MockCacheRepository) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	args := m.Called(ctx, key, value, ttl)
	return args.Error(0)
}

func (m *MockCacheRepository) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *MockCacheRepository) SetToken(ctx context.Context, token *oauth2.Token) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockCacheRepository) GetToken(ctx context.Context) (*oauth2.Token, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*oauth2.Token), args.Error(1)
}

func (m *MockCacheRepository) ClearToken(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockCacheRepository) GetUserTracks(ctx context.Context) ([]spotifyAPI.SavedTrack, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]spotifyAPI.SavedTrack), args.Error(1)
}

func (m *MockCacheRepository) SetUserTracks(ctx context.Context, tracks []spotifyAPI.SavedTrack, ttl time.Duration) error {
	args := m.Called(ctx, tracks, ttl)
	return args.Error(0)
}

func (m *MockCacheRepository) InvalidateUserTracks(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockCacheRepository) GetPlaylistTracks(ctx context.Context, id spotifyAPI.ID) ([]spotifyAPI.PlaylistTrack, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]spotifyAPI.PlaylistTrack), args.Error(1)
}

func (m *MockCacheRepository) SetPlaylistTracks(ctx context.Context, id spotifyAPI.ID, tracks []spotifyAPI.PlaylistTrack, ttl time.Duration) error {
	args := m.Called(ctx, id, tracks, ttl)
	return args.Error(0)
}

func (m *MockCacheRepository) InvalidatePlaylistTracks(ctx context.Context, id spotifyAPI.ID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

var _ shared.CacheRepository = (*MockCacheRepository)(nil)
