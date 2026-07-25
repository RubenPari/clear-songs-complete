// Package track implements application-layer use cases for track management.
// It orchestrates Spotify API calls, caching strategies, and AI-based genre
// resolution to provide track summaries, deletion operations, and artist-based
// filtering.
package track

import (
	"context"
	"sort"

	"github.com/RubenPari/clear-songs/internal/domain/shared"
	"github.com/RubenPari/clear-songs/internal/domain/track"
	"go.uber.org/zap"
)

// GetTrackSummaryUseCase handles the business logic for getting track summaries
type GetTrackSummaryUseCase struct {
	spotifyRepo shared.SpotifyRepository
	cacheRepo   shared.CacheRepository
	aiRepo      shared.AIRepository
}

// NewGetTrackSummaryUseCase creates a use case that builds per-artist track
// summaries with genre resolution. It depends on SpotifyRepository for track
// and artist data, CacheRepository for caching, and AIRepository for fallback
// genre resolution when Spotify metadata is insufficient.
func NewGetTrackSummaryUseCase(
	spotifyRepo shared.SpotifyRepository,
	cacheRepo shared.CacheRepository,
	aiRepo shared.AIRepository,
) *GetTrackSummaryUseCase {
	return &GetTrackSummaryUseCase{
		spotifyRepo: spotifyRepo,
		cacheRepo:   cacheRepo,
		aiRepo:      aiRepo,
	}
}

// Execute builds a summary of the user's saved tracks grouped by primary artist.
// It checks the cache first, and on a miss it fetches all tracks, calculates
// per-artist summaries (with genre resolution via Spotify or AI fallback),
// sorts by track count descending, and caches the result.
func (uc *GetTrackSummaryUseCase) Execute(ctx context.Context, min, max int, genre string) ([]track.ArtistSummary, error) {
	cacheKey := buildTrackSummaryCacheKey(min, max, genre)
	if cached, found := uc.getCachedSummary(ctx, cacheKey); found {
		zap.L().Info("track summary cache hit", zap.String("cache_key", cacheKey))
		return cached, nil
	}

	zap.L().Info("track summary cache miss", zap.String("cache_key", cacheKey))

	tracks, err := uc.getUserTracks(ctx)
	if err != nil {
		return nil, err
	}

	summary := uc.calculateSummary(ctx, tracks, min, max, genre)
	sort.Slice(summary, func(i, j int) bool {
		return summary[i].Count > summary[j].Count
	})

	uc.cacheSummary(ctx, cacheKey, summary)

	return summary, nil
}
