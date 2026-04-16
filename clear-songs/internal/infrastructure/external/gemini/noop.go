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

// NewNoOpAIRepository creates a new no-op AI repository
func NewNoOpAIRepository() *NoOpAIRepository {
	return &NoOpAIRepository{}
}

// ResolveArtistGenre always returns empty string (no-op)
func (n *NoOpAIRepository) ResolveArtistGenre(ctx context.Context, artistName string) (string, error) {
	noOpGenreWarn.Do(func() {
		log.Printf("[genre] GEMINI_API_KEY not set or Gemini init failed — AI fallback is a no-op (set GEMINI_API_KEY in the container env, not only .env file if /app/.env is missing)")
	})
	return "", nil
}

// Ensure NoOpAIRepository implements AIRepository
var _ shared.AIRepository = (*NoOpAIRepository)(nil)
