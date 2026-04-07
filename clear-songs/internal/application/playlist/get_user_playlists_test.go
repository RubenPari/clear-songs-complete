package playlist

import (
	"context"
	"testing"

	"github.com/RubenPari/clear-songs/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	spotifyAPI "github.com/zmb3/spotify"
)

func TestGetUserPlaylistsUseCase_Execute(t *testing.T) {
	mockSpotifyRepo := new(mocks.MockSpotifyRepository)
	mockCacheRepo := new(mocks.MockCacheRepository)
	
	useCase := NewGetUserPlaylistsUseCase(mockSpotifyRepo, mockCacheRepo)
	ctx := context.Background()

	t.Run("Success - should return playlists from API when cache miss", func(t *testing.T) {
		playlists := []spotifyAPI.SimplePlaylist{{Name: "Playlist 1"}}

		// Match any context and the specific key
		mockCacheRepo.On("Get", mock.Anything, "user_playlists", mock.Anything).Return(false, nil)
		mockSpotifyRepo.On("GetAllUserPlaylists", mock.Anything).Return(playlists, nil)
		mockCacheRepo.On("Set", mock.Anything, "user_playlists", playlists, mock.Anything).Return(nil)

		// Execute
		result, err := useCase.Execute(ctx)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, playlists, result)
	})
}
