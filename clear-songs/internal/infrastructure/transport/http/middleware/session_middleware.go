package middleware

import (
	"github.com/RubenPari/clear-songs/internal/domain/shared"
	"github.com/RubenPari/clear-songs/internal/infrastructure/logging"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Session middleware.
func SessionMiddleware(
	spotifyRepo shared.SpotifyRepository,
	cacheRepo shared.CacheRepository,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// Retrieve OAuth token from cache
		token, err := cacheRepo.GetToken(ctx)
		if err != nil {
			logging.LoggerFromGinContext(c).Error("failed to retrieve token from cache", zap.Error(err))
		}

		// If token exists, user is authenticated
		if token != nil {
			// Configure the Spotify repository with the user's token
			if err := spotifyRepo.SetAccessToken(token); err != nil {
				logging.LoggerFromGinContext(c).Error("failed to set access token", zap.Error(err))
			} else {
				// Store repositories in context for use by handlers and other middlewares
				c.Set("spotifyRepository", spotifyRepo)
				c.Set("cacheRepository", cacheRepo)
			}
		} else {
			// Log when token is not found (for debugging)
			// Only log for non-auth endpoints to avoid spam
			if c.Request.URL.Path != "/auth/is-auth" {
				logging.LoggerFromGinContext(c).Debug("no token found in cache")
			}
		}

		// Continue to next middleware or handler
		c.Next()
	}
}
