package track

import (
	"context"
	"testing"

	"github.com/RubenPari/clear-songs/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	spotifyAPI "github.com/zmb3/spotify"
)

func TestDeleteTracksByArtistUseCase_Execute(t *testing.T) {
	mockSpotifyRepo := new(mocks.MockSpotifyRepository)
	mockCacheRepo := new(mocks.MockCacheRepository)
	
	useCase := NewDeleteTracksByArtistUseCase(mockSpotifyRepo, mockCacheRepo)
	ctx := context.Background()
	artistID := spotifyAPI.ID("artist_1")

	t.Run("Success - should delete tracks and invalidate cache", func(t *testing.T) {
		tracks := []spotifyAPI.SavedTrack{
			{
				FullTrack: spotifyAPI.FullTrack{
					SimpleTrack: spotifyAPI.SimpleTrack{
						ID: "track_1",
						Artists: []spotifyAPI.SimpleArtist{{ID: "artist_1"}},
					},
				},
			},
		}
		trackIDs := []spotifyAPI.ID{"track_1"}

		mockCacheRepo.On("GetUserTracks", mock.Anything).Return(tracks, nil)
		mockSpotifyRepo.On("GetTrackIDsByArtist", mock.Anything, mock.Anything, mock.Anything).Return(trackIDs, nil)
		mockSpotifyRepo.On("DeleteTracksFromLibrary", mock.Anything, mock.Anything).Return(nil)
		mockCacheRepo.On("InvalidateUserTracks", mock.Anything).Return(nil)

		err := useCase.Execute(ctx, artistID)
		assert.NoError(t, err)
	})

	t.Run("Error - Spotify API failure should return error", func(t *testing.T) {
		tracks := []spotifyAPI.SavedTrack{{}}
		trackIDs := []spotifyAPI.ID{"track_id"}
		
		mockCacheRepo.ExpectedCalls = nil
		mockSpotifyRepo.ExpectedCalls = nil

		mockCacheRepo.On("GetUserTracks", mock.Anything).Return(tracks, nil)
		mockSpotifyRepo.On("GetTrackIDsByArtist", mock.Anything, mock.Anything, mock.Anything).Return(trackIDs, nil)
		// Force the error on the specific call
		mockSpotifyRepo.On("DeleteTracksFromLibrary", mock.Anything, mock.Anything).Return(assert.AnError)

		err := useCase.Execute(ctx, artistID)

		assert.Error(t, err, "Should return error when Spotify API fails")
	})
}
