package middleware

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/RubenPari/clear-songs/internal/domain/shared"
	"github.com/gin-gonic/gin"
	spotifyAPI "github.com/zmb3/spotify"
)

// CacheInvalidationMiddleware invalidates cache after successful write operations.
func CacheInvalidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if !isWriteMethod(c.Request.Method) {
			return
		}

		if status := c.Writer.Status(); status < 200 || status >= 300 {
			return
		}

		cacheRepo, ok := cacheRepositoryFromContext(c)
		if !ok {
			return
		}

		if err := invalidateByRequest(c, cacheRepo); err != nil {
			log.Printf("warning: cache invalidation failed for %s %s: %v", c.Request.Method, c.Request.URL.Path, err)
		}
	}
}

func cacheRepositoryFromContext(c *gin.Context) (shared.CacheRepository, bool) {
	repo, exists := c.Get("cacheRepository")
	if !exists {
		return nil, false
	}

	cacheRepo, ok := repo.(shared.CacheRepository)
	if !ok || cacheRepo == nil {
		return nil, false
	}

	return cacheRepo, true
}

func isWriteMethod(method string) bool {
	switch method {
	case "DELETE", "POST", "PUT", "PATCH":
		return true
	default:
		return false
	}
}

func invalidateByRequest(c *gin.Context, cacheRepo shared.CacheRepository) error {
	ctx := c.Request.Context()
	path := c.Request.URL.Path

	switch {
	case strings.HasPrefix(path, "/track/"):
		if err := cacheRepo.InvalidateUserTracks(ctx); err != nil {
			return fmt.Errorf("invalidate user tracks: %w", err)
		}
		return nil
	case strings.HasPrefix(path, "/playlist/"):
		return invalidatePlaylistCaches(ctx, c, cacheRepo, path)
	case strings.HasPrefix(path, "/album/"):
		if err := cacheRepo.InvalidateUserTracks(ctx); err != nil {
			return fmt.Errorf("invalidate user tracks: %w", err)
		}
		return nil
	default:
		if err := cacheRepo.InvalidateUserTracks(ctx); err != nil {
			return fmt.Errorf("invalidate fallback user tracks: %w", err)
		}
		return nil
	}
}

func invalidatePlaylistCaches(ctx context.Context, c *gin.Context, cacheRepo shared.CacheRepository, path string) error {
	var errs []error

	if playlistID := c.Query("id"); playlistID != "" {
		if err := cacheRepo.InvalidatePlaylistTracks(ctx, spotifyAPI.ID(playlistID)); err != nil {
			errs = append(errs, fmt.Errorf("invalidate playlist tracks: %w", err))
		}
	}

	if strings.Contains(path, "library") {
		if err := cacheRepo.InvalidateUserTracks(ctx); err != nil {
			errs = append(errs, fmt.Errorf("invalidate user tracks: %w", err))
		}
	}

	return errors.Join(errs...)
}
