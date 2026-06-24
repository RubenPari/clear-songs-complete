package playlist

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/RubenPari/clear-songs/internal/domain/shared"
	"github.com/RubenPari/clear-songs/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	spotifyAPI "github.com/zmb3/spotify"
)

func TestExtractPlaylistTrackIDs(t *testing.T) {
	t.Run("extracts IDs from tracks", func(t *testing.T) {
		tracks := []spotifyAPI.PlaylistTrack{
			{Track: spotifyAPI.FullTrack{SimpleTrack: spotifyAPI.SimpleTrack{ID: "id1"}}},
			{Track: spotifyAPI.FullTrack{SimpleTrack: spotifyAPI.SimpleTrack{ID: "id2"}}},
		}

		ids := extractPlaylistTrackIDs(tracks)

		assert.Equal(t, []spotifyAPI.ID{"id1", "id2"}, ids)
	})

	t.Run("empty input returns empty slice", func(t *testing.T) {
		ids := extractPlaylistTrackIDs(nil)

		assert.Empty(t, ids)
		assert.NotNil(t, ids)
	})
}

func TestFetchPlaylistTracks(t *testing.T) {
	ctx := context.Background()
	playlistID := spotifyAPI.ID("playlist1")

	t.Run("returns cached tracks on cache hit", func(t *testing.T) {
		mockSpotifyRepo := new(mocks.MockSpotifyRepository)
		mockCacheRepo := new(mocks.MockCacheRepository)

		cached := []spotifyAPI.PlaylistTrack{
			{Track: spotifyAPI.FullTrack{SimpleTrack: spotifyAPI.SimpleTrack{ID: "t1"}}},
		}
		mockCacheRepo.On("GetPlaylistTracks", ctx, playlistID).Return(cached, nil)

		tracks, err := fetchPlaylistTracks(ctx, mockSpotifyRepo, mockCacheRepo, playlistID)

		assert.NoError(t, err)
		assert.Equal(t, cached, tracks)
		mockCacheRepo.AssertExpectations(t)
		mockSpotifyRepo.AssertNotCalled(t, "GetAllPlaylistTracks")
	})

	t.Run("fetches from API and caches on cache miss", func(t *testing.T) {
		mockSpotifyRepo := new(mocks.MockSpotifyRepository)
		mockCacheRepo := new(mocks.MockCacheRepository)

		fetched := []spotifyAPI.PlaylistTrack{
			{Track: spotifyAPI.FullTrack{SimpleTrack: spotifyAPI.SimpleTrack{ID: "t2"}}},
		}
		mockCacheRepo.On("GetPlaylistTracks", ctx, playlistID).Return(nil, nil)
		mockSpotifyRepo.On("GetAllPlaylistTracks", ctx, playlistID).Return(fetched, nil)
		mockCacheRepo.On("SetPlaylistTracks", ctx, playlistID, fetched, mock.MatchedBy(func(ttl time.Duration) bool {
			return ttl == 5*time.Minute
		})).Return(nil)

		tracks, err := fetchPlaylistTracks(ctx, mockSpotifyRepo, mockCacheRepo, playlistID)

		assert.NoError(t, err)
		assert.Equal(t, fetched, tracks)
		mockCacheRepo.AssertExpectations(t)
		mockSpotifyRepo.AssertExpectations(t)
	})

	t.Run("returns API error without caching", func(t *testing.T) {
		mockSpotifyRepo := new(mocks.MockSpotifyRepository)
		mockCacheRepo := new(mocks.MockCacheRepository)

		apiErr := errors.New("spotify failure")
		mockCacheRepo.On("GetPlaylistTracks", ctx, playlistID).Return(nil, nil)
		mockSpotifyRepo.On("GetAllPlaylistTracks", ctx, playlistID).Return(nil, apiErr)

		tracks, err := fetchPlaylistTracks(ctx, mockSpotifyRepo, mockCacheRepo, playlistID)

		assert.ErrorIs(t, err, apiErr)
		assert.Nil(t, tracks)
		mockCacheRepo.AssertNotCalled(t, "SetPlaylistTracks")
	})

	t.Run("uses shared cache key format for playlist tracks", func(t *testing.T) {
		mockSpotifyRepo := new(mocks.MockSpotifyRepository)
		mockCacheRepo := new(mocks.MockCacheRepository)

		fetched := []spotifyAPI.PlaylistTrack{
			{Track: spotifyAPI.FullTrack{SimpleTrack: spotifyAPI.SimpleTrack{ID: "t3"}}},
		}
		mockCacheRepo.On("GetPlaylistTracks", ctx, playlistID).Return(nil, nil)
		mockSpotifyRepo.On("GetAllPlaylistTracks", ctx, playlistID).Return(fetched, nil)
		mockCacheRepo.On("SetPlaylistTracks", ctx, playlistID, fetched, mock.Anything).Return(nil)

		expectedKey := shared.CacheKeyPlaylistTracks(string(playlistID))
		assert.Equal(t, "tracksPlaylistplaylist1", expectedKey)

		_, err := fetchPlaylistTracks(ctx, mockSpotifyRepo, mockCacheRepo, playlistID)

		assert.NoError(t, err)
	})
}
