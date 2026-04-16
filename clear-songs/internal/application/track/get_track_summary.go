package track

import (
	"context"
	"log"
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

// Creates get track summary use case.
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

// Execute.
func (uc *GetTrackSummaryUseCase) Execute(ctx context.Context, min, max int, genre string) ([]track.ArtistSummary, error) {
	cacheKey := buildTrackSummaryCacheKey(min, max, genre)
	if cached, found := uc.getCachedSummary(ctx, cacheKey); found {
		log.Printf("[track/summary] Redis cache HIT key=%q — genre resolution and [genre] AI logs are skipped. POST /track/library-cache/invalidate then GET /track/summary again to recompute.", cacheKey)
		return cached, nil
	}

	log.Printf("[track/summary] cache MISS key=%q — computing (per-artist [genre] logs may appear)", cacheKey)

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
