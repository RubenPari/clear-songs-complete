package gemini

import (
	"context"

	"github.com/RubenPari/clear-songs/internal/domain/shared"
)

// NoOpAIRepository is a no-op implementation of AIRepository
// used when Gemini API key is not configured
type NoOpAIRepository struct{}

// NewNoOpAIRepository creates a new no-op AI repository
func NewNoOpAIRepository() *NoOpAIRepository {
	return &NoOpAIRepository{}
}

// ResolveArtistGenre always returns empty string (no-op)
func (n *NoOpAIRepository) ResolveArtistGenre(ctx context.Context, artistName string) (string, error) {
	return "", nil
}

// Ensure NoOpAIRepository implements AIRepository
var _ shared.AIRepository = (*NoOpAIRepository)(nil)
