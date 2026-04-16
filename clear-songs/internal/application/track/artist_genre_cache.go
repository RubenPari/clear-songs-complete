package track

import (
	"context"
	"os"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
)

const artistAIGenreKeyPrefix = "artist_ai_genre:"

// Artist aigenre cache key.
func artistAIGenreCacheKey(artistKey string) string {
	return artistAIGenreKeyPrefix + artistKey
}

// Artist aigenre cache ttl.
func artistAIGenreCacheTTL() time.Duration {
	const defaultSec = 7 * 24 * 3600
	s := strings.TrimSpace(os.Getenv("ARTIST_AI_GENRE_CACHE_TTL_SEC"))
	if s == "" {
		return time.Duration(defaultSec) * time.Second
	}
	sec, err := strconv.Atoi(s)
	if err != nil || sec < 60 {
		return time.Duration(defaultSec) * time.Second
	}
	return time.Duration(sec) * time.Second
}

// Fetches cached artist canonical genre.
func (uc *GetTrackSummaryUseCase) getCachedArtistCanonicalGenre(ctx context.Context, artistKey string) (string, bool) {
	if uc.cacheRepo == nil {
		return "", false
	}
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
	if uc.cacheRepo == nil || strings.TrimSpace(canonical) == "" {
		return
	}
	if err := uc.cacheRepo.Set(ctx, artistAIGenreCacheKey(artistKey), canonical, artistAIGenreCacheTTL()); err != nil {
		zap.L().Warn("artist genre cache write failed", zap.String("artist_key", artistKey), zap.Error(err))
	}
}
