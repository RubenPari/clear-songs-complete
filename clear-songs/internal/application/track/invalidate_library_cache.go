package track

import (
	"context"
	"fmt"

	"github.com/RubenPari/clear-songs/internal/domain/shared"
)

// InvalidateLibraryCacheUseCase handles cache invalidation for the user's library.
type InvalidateLibraryCacheUseCase struct {
	cacheRepo shared.CacheRepository
}

// NewInvalidateLibraryCacheUseCase creates a new InvalidateLibraryCacheUseCase.
func NewInvalidateLibraryCacheUseCase(cacheRepo shared.CacheRepository) *InvalidateLibraryCacheUseCase {
	return &InvalidateLibraryCacheUseCase{cacheRepo: cacheRepo}
}

// Execute invalidates the cached user tracks and summary entries.
// It returns the message that should be sent back to the client.
func (uc *InvalidateLibraryCacheUseCase) Execute(ctx context.Context) (string, error) {
	if err := uc.cacheRepo.InvalidateUserTracks(ctx); err != nil {
		return "", fmt.Errorf("invalidate user tracks cache: %w", err)
	}
	return "Library cache invalidated", nil
}
