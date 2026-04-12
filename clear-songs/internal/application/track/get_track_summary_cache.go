package track

import (
	"context"
	"fmt"
	"time"

	domainTrack "github.com/RubenPari/clear-songs/internal/domain/track"
)

func buildTrackSummaryCacheKey(min, max int, genre string) string {
	if min > 0 || max > 0 {
		return fmt.Sprintf("track_summary_%d_%d_%s", min, max, genre)
	}

	return "track_summary"
}

func (uc *GetTrackSummaryUseCase) getCachedSummary(ctx context.Context, cacheKey string) ([]domainTrack.ArtistSummary, bool) {
	if uc.cacheRepo == nil {
		return nil, false
	}

	var cached []domainTrack.ArtistSummary
	found, _ := uc.cacheRepo.Get(ctx, cacheKey, &cached)
	if !found {
		return nil, false
	}

	return cached, true
}

func (uc *GetTrackSummaryUseCase) cacheSummary(ctx context.Context, cacheKey string, summary []domainTrack.ArtistSummary) {
	if uc.cacheRepo == nil {
		return
	}

	_ = uc.cacheRepo.Set(ctx, cacheKey, summary, 5*time.Minute)
}
