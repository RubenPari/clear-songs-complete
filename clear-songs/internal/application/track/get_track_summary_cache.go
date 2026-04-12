package track

import (
	"context"
	"fmt"
	"log"
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
	found, err := uc.cacheRepo.Get(ctx, cacheKey, &cached)
	if err != nil {
		log.Printf("warning: failed to read summary cache for key %s: %v", cacheKey, err)
		return nil, false
	}
	if !found {
		return nil, false
	}

	return cached, true
}

func (uc *GetTrackSummaryUseCase) cacheSummary(ctx context.Context, cacheKey string, summary []domainTrack.ArtistSummary) {
	if uc.cacheRepo == nil {
		return
	}

	if err := uc.cacheRepo.Set(ctx, cacheKey, summary, 5*time.Minute); err != nil {
		log.Printf("warning: failed to write summary cache for key %s: %v", cacheKey, err)
	}
}
