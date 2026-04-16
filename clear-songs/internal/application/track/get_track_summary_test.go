package track

import (
	"context"
	"testing"

	"github.com/RubenPari/clear-songs/internal/domain/shared"
	"github.com/RubenPari/clear-songs/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/mock"
	spotifyAPI "github.com/zmb3/spotify"
)

func TestBuildTrackSummaryCacheKey(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "track_summary_0_0", buildTrackSummaryCacheKey(0, 0, ""))
	assert.Equal(t, "track_summary_0_0_Rock", buildTrackSummaryCacheKey(0, 0, "Rock"))
	assert.Equal(t, "track_summary_5_50", buildTrackSummaryCacheKey(5, 50, ""))
	assert.Equal(t, "track_summary_5_50_Hip Hop", buildTrackSummaryCacheKey(5, 50, "Hip Hop"))
}

// TestGetTrackSummaryUseCase_Execute tests the Execute method of GetTrackSummaryUseCase
// It verifies that the method correctly handles both cache hits and misses
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
		mockCacheRepo.On("Get", ctx, "track_summary_0_0", mock.Anything).Return(false, nil)
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

		mockCacheRepo.On("Set", ctx, "track_summary_0_0", mock.Anything, mock.Anything).Return(nil)

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

func TestGetTrackSummaryUseCase_Execute_AI_batch_empty_spotify_genres(t *testing.T) {
	mockSpotifyRepo := new(mocks.MockSpotifyRepository)
	mockCacheRepo := new(mocks.MockCacheRepository)
	mockAIRepo := new(mocks.MockAIRepository)
	useCase := NewGetTrackSummaryUseCase(mockSpotifyRepo, mockCacheRepo, mockAIRepo)
	ctx := context.Background()

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
	}

	mockCacheRepo.On("Get", ctx, "track_summary_0_0", mock.Anything).Return(false, nil).Once()
	mockCacheRepo.On("GetUserTracks", ctx).Return(nil, nil)
	mockSpotifyRepo.On("GetAllUserTracks", ctx).Return(tracks, nil)
	mockCacheRepo.On("SetUserTracks", ctx, tracks, mock.Anything).Return(nil)

	mockSpotifyRepo.On("GetArtists", ctx, mock.MatchedBy(func(ids []spotifyAPI.ID) bool {
		return len(ids) == 1 && string(ids[0]) == "1"
	})).Return([]*spotifyAPI.FullArtist{
		{
			SimpleArtist: spotifyAPI.SimpleArtist{ID: "1", Name: "Artist 1"},
			Genres:       nil,
		},
	}, nil)

	mockCacheRepo.On("Get", ctx, "artist_ai_genre:1", mock.Anything).Return(false, nil).Once()
	mockAIRepo.On("ResolveArtistGenres", ctx, mock.MatchedBy(func(lookups []shared.AIArtistLookup) bool {
		return len(lookups) == 1 && lookups[0].Key == "1" && lookups[0].Name == "Artist 1"
	})).Return(map[string]string{"1": "hip hop"}, nil).Once()
	mockCacheRepo.On("Set", ctx, "artist_ai_genre:1", "Hip Hop", mock.Anything).Return(nil).Once()
	mockCacheRepo.On("Set", ctx, "track_summary_0_0", mock.Anything, mock.Anything).Return(nil)

	result, err := useCase.Execute(ctx, 0, 0, "")

	assert.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, "Artist 1", result[0].Name)
	assert.Equal(t, "Hip Hop", result[0].Genre)
}
