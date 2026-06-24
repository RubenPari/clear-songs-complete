package track

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/RubenPari/clear-songs/internal/domain/shared"
	"go.uber.org/zap"
)

var (
	artistAIGenreTTL     time.Duration
	artistAIGenreTTLOnce sync.Once
)

const artistAIGenreKeyPrefix = "artist_ai_genre:"

// Artist aigenre cache key.
func artistAIGenreCacheKey(artistKey string) string {
	return fmt.Sprintf(shared.CacheKeyArtistAIGenreFmt, artistKey)
}

// Artist aigenre cache ttl reads the environment once and caches the result.
func artistAIGenreCacheTTL() time.Duration {
	artistAIGenreTTLOnce.Do(func() {
		const defaultSec = 7 * 24 * 3600
		s := strings.TrimSpace(os.Getenv("ARTIST_AI_GENRE_CACHE_TTL_SEC"))
		if s == "" {
			artistAIGenreTTL = time.Duration(defaultSec) * time.Second
			return
		}
		sec, err := strconv.Atoi(s)
		if err != nil || sec < 60 {
			artistAIGenreTTL = time.Duration(defaultSec) * time.Second
			return
		}
		artistAIGenreTTL = time.Duration(sec) * time.Second
	})
	return artistAIGenreTTL
}

// Fetches cached artist canonical genre.
func (uc *GetTrackSummaryUseCase) getCachedArtistCanonicalGenre(ctx context.Context, artistKey string) (string, bool) {
	var s string
	found, err := uc.cacheRepo.Get(ctx, artistAIGenreCacheKey(artistKey), &s)
	if err != nil {
		zap.L().Warn("artist genre cache read failed", zap.String("artist_key", artistKey), zap.Error(err))
		return "", false
	}
	if !found || strings.TrimSpace(s) == "" {
		return "", false
	}
	return s, true
}

// Sets cached artist canonical genre.
func (uc *GetTrackSummaryUseCase) setCachedArtistCanonicalGenre(ctx context.Context, artistKey, canonical string) {
	if strings.TrimSpace(canonical) == "" {
		return
	}
	if err := uc.cacheRepo.Set(ctx, artistAIGenreCacheKey(artistKey), canonical, artistAIGenreCacheTTL()); err != nil {
		zap.L().Warn("artist genre cache write failed", zap.String("artist_key", artistKey), zap.Error(err))
	}
}
