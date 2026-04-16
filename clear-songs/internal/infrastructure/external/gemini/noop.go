package gemini

import (
	"context"
	"log"
	"sync"

	"github.com/RubenPari/clear-songs/internal/domain/shared"
)

// NoOpAIRepository is a no-op implementation of AIRepository
// used when Gemini API key is not configured
type NoOpAIRepository struct{}

var noOpGenreWarn sync.Once

// Creates no op airepository.
func NewNoOpAIRepository() *NoOpAIRepository {
	return &NoOpAIRepository{}
}

// Resolves artist genres.
func (n *NoOpAIRepository) ResolveArtistGenres(ctx context.Context, lookups []shared.AIArtistLookup) (map[string]string, error) {
	if len(lookups) > 0 {
		noOpGenreWarn.Do(func() {
			log.Printf("[genre] GEMINI_API_KEY not set or Gemini init failed — AI fallback is a no-op (set GEMINI_API_KEY in the container env, not only .env file if /app/.env is missing)")
		})
	}
	out := make(map[string]string, len(lookups))
	for _, l := range lookups {
		out[l.Key] = ""
	}
	return out, nil
}

// Ensure NoOpAIRepository implements AIRepository
var _ shared.AIRepository = (*NoOpAIRepository)(nil)
