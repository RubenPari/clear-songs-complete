package track

import (
	"context"
	"fmt"
	"time"

	"github.com/RubenPari/clear-songs/internal/domain/shared"
	domainTrack "github.com/RubenPari/clear-songs/internal/domain/track"
	"go.uber.org/zap"
)

// buildTrackSummaryCacheKey constructs a cache key for the track summary based
// on the range filter (min, max) and optional genre filter.
func buildTrackSummaryCacheKey(min, max int, genre string) string {
	if genre == "" {
		return fmt.Sprintf(shared.CacheKeyTrackSummaryFmt, min, max)
	}
	return fmt.Sprintf(shared.CacheKeyTrackSummaryGenreFmt, min, max, genre)
}

// getCachedSummary attempts to read a previously cached track summary from Redis.
// Returns the summary and true on cache hit, or nil and false on miss or error.
func (uc *GetTrackSummaryUseCase) getCachedSummary(ctx context.Context, cacheKey string) ([]domainTrack.ArtistSummary, bool) {
	var cached []domainTrack.ArtistSummary
	found, err := uc.cacheRepo.Get(ctx, cacheKey, &cached)
	if err != nil {
		zap.L().Warn("failed to read summary cache", zap.String("cache_key", cacheKey), zap.Error(err))
		return nil, false
	}
	if !found {
		return nil, false
	}

	return cached, true
}

// cacheSummary writes the track summary to Redis with a 5-minute TTL. Errors are
// logged but do not propagate, as caching is a performance optimisation.
func (uc *GetTrackSummaryUseCase) cacheSummary(ctx context.Context, cacheKey string, summary []domainTrack.ArtistSummary) {
	if err := uc.cacheRepo.Set(ctx, cacheKey, summary, 5*time.Minute); err != nil {
		zap.L().Warn("failed to write summary cache", zap.String("cache_key", cacheKey), zap.Error(err))
	}
}
