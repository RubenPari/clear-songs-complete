package track

import (
	"context"
	"sort"

	"github.com/RubenPari/clear-songs/internal/domain/shared"
	"github.com/RubenPari/clear-songs/internal/domain/track"
)

// GetTrackSummaryUseCase handles the business logic for getting track summaries
type GetTrackSummaryUseCase struct {
	spotifyRepo shared.SpotifyRepository
	cacheRepo   shared.CacheRepository
	aiRepo      shared.AIRepository
}

// NewGetTrackSummaryUseCase creates a new GetTrackSummaryUseCase
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

// Execute retrieves track summary grouped by artist, optionally filtered by range and genre
func (uc *GetTrackSummaryUseCase) Execute(ctx context.Context, min, max int, genre string) ([]track.ArtistSummary, error) {
	cacheKey := buildTrackSummaryCacheKey(min, max, genre)
	if cached, found := uc.getCachedSummary(ctx, cacheKey); found {
		return cached, nil
	}

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
