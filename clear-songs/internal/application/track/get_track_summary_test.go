package track

import (
	"context"
	"testing"

	"github.com/RubenPari/clear-songs/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	spotifyAPI "github.com/zmb3/spotify"
)

func TestGetTrackSummaryUseCase_Execute(t *testing.T) {
	// Setup mocks
	mockSpotifyRepo := new(mocks.MockSpotifyRepository)
	mockCacheRepo := new(mocks.MockCacheRepository)
	mockAIRepo := new(mocks.MockAIRepository)

	useCase := NewGetTrackSummaryUseCase(mockSpotifyRepo, mockCacheRepo, mockAIRepo)
	ctx := context.Background()

	t.Run("Success - should return grouped summary with resolved genres", func(t *testing.T) {
		// Mock data
		tracks := []spotifyAPI.SavedTrack{
			{
				FullTrack: spotifyAPI.FullTrack{
					SimpleTrack: spotifyAPI.SimpleTrack{
						Name: "Song A",
						Artists: []spotifyAPI.SimpleArtist{
							{Name: "Artist 1", ID: "1"},
						},
					},
				},
			},
			{
				FullTrack: spotifyAPI.FullTrack{
					SimpleTrack: spotifyAPI.SimpleTrack{
						Name: "Song B",
						Artists: []spotifyAPI.SimpleArtist{
							{Name: "Artist 1", ID: "1"},
						},
					},
				},
			},
			{
				FullTrack: spotifyAPI.FullTrack{
					SimpleTrack: spotifyAPI.SimpleTrack{
						Name: "Song C",
						Artists: []spotifyAPI.SimpleArtist{
							{Name: "Artist 2", ID: "2"},
						},
					},
				},
			},
		}

		// Configure mock expectations
		mockCacheRepo.On("Get", ctx, "track_summary", mock.Anything).Return(false, nil)
		mockCacheRepo.On("GetUserTracks", ctx).Return(nil, nil)
		mockSpotifyRepo.On("GetAllUserTracks", ctx).Return(tracks, nil)
		mockCacheRepo.On("SetUserTracks", ctx, tracks, mock.Anything).Return(nil)

		// Mock batch artist fetch with genres
		mockSpotifyRepo.On("GetArtists", ctx, mock.MatchedBy(func(ids []spotifyAPI.ID) bool {
			return len(ids) == 2
		})).Return([]*spotifyAPI.FullArtist{
			{
				SimpleArtist: spotifyAPI.SimpleArtist{ID: "1", Name: "Artist 1"},
				Genres:       []string{"gangster rap", "east coast hip hop"},
			},
			{
				SimpleArtist: spotifyAPI.SimpleArtist{ID: "2", Name: "Artist 2"},
				Genres:       []string{"classic rock"},
			},
		}, nil)

		mockCacheRepo.On("Set", ctx, "track_summary", mock.Anything, mock.Anything).Return(nil)

		// Execute
		result, err := useCase.Execute(ctx, 0, 0, "")

		// Assertions
		assert.NoError(t, err)
		assert.Len(t, result, 2)

		// Check Artist 1 (should have 2 tracks, genre resolved to Hip Hop)
		assert.Equal(t, "Artist 1", result[0].Name)
		assert.Equal(t, 2, result[0].Count)
		assert.Equal(t, "Hip Hop", result[0].Genre)

		// Check Artist 2 (should have 1 track, genre resolved to Rock)
		assert.Equal(t, "Artist 2", result[1].Name)
		assert.Equal(t, 1, result[1].Count)
		assert.Equal(t, "Rock", result[1].Genre)
	})
}
