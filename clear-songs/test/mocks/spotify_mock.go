package mocks

import (
	"context"
	"time"

	"github.com/RubenPari/clear-songs/internal/domain/shared"
	"github.com/stretchr/testify/mock"
	spotifyAPI "github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

// MockSpotifyRepository is a mock implementation of SpotifyRepository
type MockSpotifyRepository struct {
	mock.Mock
}

func (m *MockSpotifyRepository) GetCurrentUser(ctx context.Context) (*spotifyAPI.PrivateUser, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*spotifyAPI.PrivateUser), args.Error(1)
}

func (m *MockSpotifyRepository) GetAllUserTracks(ctx context.Context) ([]spotifyAPI.SavedTrack, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]spotifyAPI.SavedTrack), args.Error(1)
}

func (m *MockSpotifyRepository) GetArtist(ctx context.Context, id spotifyAPI.ID) (*spotifyAPI.FullArtist, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*spotifyAPI.FullArtist), args.Error(1)
}

func (m *MockSpotifyRepository) GetArtists(ctx context.Context, ids []spotifyAPI.ID) ([]*spotifyAPI.FullArtist, error) {
	args := m.Called(ctx, ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*spotifyAPI.FullArtist), args.Error(1)
}

func (m *MockSpotifyRepository) GetTrack(ctx context.Context, id spotifyAPI.ID) (*spotifyAPI.FullTrack, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*spotifyAPI.FullTrack), args.Error(1)
}

func (m *MockSpotifyRepository) GetTrackIDsByArtist(ctx context.Context, id spotifyAPI.ID, tracks []spotifyAPI.SavedTrack) ([]spotifyAPI.ID, error) {
	args := m.Called(ctx, id, tracks)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]spotifyAPI.ID), args.Error(1)
}

func (m *MockSpotifyRepository) DeleteTracksFromLibrary(ctx context.Context, ids []spotifyAPI.ID) error {
	args := m.Called(ctx, ids)
	return args.Error(0)
}

func (m *MockSpotifyRepository) GetAllUserPlaylists(ctx context.Context) ([]spotifyAPI.SimplePlaylist, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]spotifyAPI.SimplePlaylist), args.Error(1)
}

// Minimal behavior for remaining methods
func (m *MockSpotifyRepository) GetUserTracks(ctx context.Context, limit, offset int) ([]spotifyAPI.SavedTrack, error) {
	return nil, nil
}
func (m *MockSpotifyRepository) GetTracksByArtist(ctx context.Context, id spotifyAPI.ID, tracks []spotifyAPI.SavedTrack) ([]spotifyAPI.SavedTrack, error) {
	return nil, nil
}
func (m *MockSpotifyRepository) GetPlaylist(ctx context.Context, id spotifyAPI.ID) (*spotifyAPI.FullPlaylist, error) {
	return nil, nil
}
func (m *MockSpotifyRepository) GetPlaylistTracks(ctx context.Context, id spotifyAPI.ID, limit, offset int) ([]spotifyAPI.PlaylistTrack, error) {
	return nil, nil
}
func (m *MockSpotifyRepository) GetAllPlaylistTracks(ctx context.Context, id spotifyAPI.ID) ([]spotifyAPI.PlaylistTrack, error) {
	return nil, nil
}
func (m *MockSpotifyRepository) DeletePlaylistTracks(ctx context.Context, id spotifyAPI.ID, ids []spotifyAPI.ID) error {
	return nil
}
func (m *MockSpotifyRepository) GetUserPlaylists(ctx context.Context, limit, offset int) ([]spotifyAPI.SimplePlaylist, error) {
	return nil, nil
}
func (m *MockSpotifyRepository) SetAccessToken(token interface{}) error {
	return nil
}

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

// Minimal behavior for remaining methods
func (m *MockCacheRepository) SetToken(ctx context.Context, token *oauth2.Token) error { return nil }
func (m *MockCacheRepository) GetToken(ctx context.Context) (*oauth2.Token, error)     { return nil, nil }
func (m *MockCacheRepository) ClearToken(ctx context.Context) error                    { return nil }
func (m *MockCacheRepository) GetPlaylistTracks(ctx context.Context, id spotifyAPI.ID) ([]spotifyAPI.PlaylistTrack, error) {
	return nil, nil
}
func (m *MockCacheRepository) SetPlaylistTracks(ctx context.Context, id spotifyAPI.ID, tracks []spotifyAPI.PlaylistTrack, ttl time.Duration) error {
	return nil
}
func (m *MockCacheRepository) InvalidatePlaylistTracks(ctx context.Context, id spotifyAPI.ID) error {
	return nil
}
func (m *MockCacheRepository) Delete(ctx context.Context, key string) error { return nil }

var _ shared.SpotifyRepository = (*MockSpotifyRepository)(nil)
var _ shared.CacheRepository = (*MockCacheRepository)(nil)
