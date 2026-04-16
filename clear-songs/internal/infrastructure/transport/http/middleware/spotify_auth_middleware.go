package middleware

import (
	"net/http"

	"github.com/RubenPari/clear-songs/internal/application/shared/dto"
	"github.com/RubenPari/clear-songs/internal/domain/shared"
	"github.com/gin-gonic/gin"
)

// Spotify auth middleware.
func SpotifyAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Spotify repository from context (set by SessionMiddleware)
		repo, exists := c.Get("spotifyRepository")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.UnauthorizedErr())
			return
		}

		spotifyRepo, ok := repo.(shared.SpotifyRepository)
		if !ok || spotifyRepo == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.UnauthorizedErr())
			return
		}

		c.Set("spotifyRepository", spotifyRepo)
		c.Next()
	}
}
