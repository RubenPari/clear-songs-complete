package mocks

import (
	"context"

	"github.com/RubenPari/clear-songs/internal/domain/shared"
	"github.com/stretchr/testify/mock"
	spotifyAPI "github.com/zmb3/spotify"
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
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]spotifyAPI.PlaylistTrack), args.Error(1)
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

var _ shared.SpotifyRepository = (*MockSpotifyRepository)(nil)
