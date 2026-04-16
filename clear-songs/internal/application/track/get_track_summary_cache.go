package track

import (
	"context"
	"fmt"
	"time"

	domainTrack "github.com/RubenPari/clear-songs/internal/domain/track"
	"go.uber.org/zap"
)

// Builds track summary cache key.
func buildTrackSummaryCacheKey(min, max int, genre string) string {
	if genre == "" {
		return fmt.Sprintf("track_summary_%d_%d", min, max)
	}
	return fmt.Sprintf("track_summary_%d_%d_%s", min, max, genre)
}

// Fetches cached summary.
func (uc *GetTrackSummaryUseCase) getCachedSummary(ctx context.Context, cacheKey string) ([]domainTrack.ArtistSummary, bool) {
	if uc.cacheRepo == nil {
		return nil, false
	}

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

// Cache summary.
func (uc *GetTrackSummaryUseCase) cacheSummary(ctx context.Context, cacheKey string, summary []domainTrack.ArtistSummary) {
	if uc.cacheRepo == nil {
		return
	}

	if err := uc.cacheRepo.Set(ctx, cacheKey, summary, 5*time.Minute); err != nil {
		zap.L().Warn("failed to write summary cache", zap.String("cache_key", cacheKey), zap.Error(err))
	}
}
