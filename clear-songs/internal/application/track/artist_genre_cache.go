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

// artistAIGenreCacheKey constructs a Redis cache key for the AI-resolved canonical
// genre of a specific artist (identified by artistKey).
func artistAIGenreCacheKey(artistKey string) string {
	return fmt.Sprintf(shared.CacheKeyArtistAIGenreFmt, artistKey)
}

// artistAIGenreCacheTTL reads the ARTIST_AI_GENRE_CACHE_TTL_SEC environment variable
// once and returns the corresponding duration. Falls back to 7 days if the variable
// is unset, invalid, or below 60 seconds.
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

// getCachedArtistCanonicalGenre attempts to read the AI-resolved canonical genre
// for an artist from Redis. Returns the genre and true on cache hit, or empty
// string and false on miss or error.
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

// setCachedArtistCanonicalGenre writes the AI-resolved canonical genre for an artist
// to Redis with the configured TTL. Empty values are not cached. Errors are logged
// but do not propagate.
func (uc *GetTrackSummaryUseCase) setCachedArtistCanonicalGenre(ctx context.Context, artistKey, canonical string) {
	if strings.TrimSpace(canonical) == "" {
		return
	}
	if err := uc.cacheRepo.Set(ctx, artistAIGenreCacheKey(artistKey), canonical, artistAIGenreCacheTTL()); err != nil {
		zap.L().Warn("artist genre cache write failed", zap.String("artist_key", artistKey), zap.Error(err))
	}
}
