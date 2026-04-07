package middleware

import (
	"strings"
	"github.com/RubenPari/clear-songs/internal/domain/shared"
	"github.com/gin-gonic/gin"
	spotifyAPI "github.com/zmb3/spotify"
)

// CacheInvalidationMiddleware automatically invalidates cache based on the endpoint called
func CacheInvalidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Execute the request first
		c.Next()

		// Only invalidate cache if the request was successful
		if c.Writer.Status() >= 200 && c.Writer.Status() < 300 {
			// Get CacheRepository from context (set by SessionMiddleware)
			repo, exists := c.Get("cacheRepository")
			if !exists {
				return
			}

			cacheRepo, ok := repo.(shared.CacheRepository)
			if !ok || cacheRepo == nil {
				return
			}

			path := c.Request.URL.Path
			method := c.Request.Method

			// Only invalidate on modification operations (DELETE, POST, PUT, PATCH)
			if method == "DELETE" || method == "POST" || method == "PUT" || method == "PATCH" {
				invalidateBasedOnEndpoint(c, cacheRepo, path)
			}
		}
	}
}

func invalidateBasedOnEndpoint(c *gin.Context, cacheRepo shared.CacheRepository, path string) {
	ctx := c.Request.Context()

	switch {
	case strings.HasPrefix(path, "/track/"):
		// Any track operation affects user data
		_ = cacheRepo.InvalidateUserTracks(ctx)

	case strings.HasPrefix(path, "/playlist/"):
		// Playlist operations
		if playlistID := c.Query("id"); playlistID != "" {
			_ = cacheRepo.InvalidatePlaylistTracks(ctx, spotifyAPI.ID(playlistID))
		}

		// If it's a playlist operation that also affects user library
		if strings.Contains(path, "all") || strings.Contains(path, "library") {
			_ = cacheRepo.InvalidateUserTracks(ctx)
		}

	case strings.HasPrefix(path, "/album/"):
		// Album operations usually affect user library
		_ = cacheRepo.InvalidateUserTracks(ctx)

	default:
		// For any other modification operation, do a full reset as a safety measure
		// Note: CacheRepository doesn't have a Reset method, so we invalidate user tracks at minimum
		_ = cacheRepo.InvalidateUserTracks(ctx)
	}
}
