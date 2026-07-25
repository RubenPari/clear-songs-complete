// Package middleware provides HTTP middleware for request processing, authentication,
// and session management. It integrates with the Gin framework and enforces security
// policies at the transport layer.
package middleware

import (
	"net/http"

	"github.com/RubenPari/clear-songs/internal/application/shared/dto"
	"github.com/RubenPari/clear-songs/internal/domain/shared"
	"github.com/gin-gonic/gin"
)

// SpotifyAuthMiddleware verifies that a Spotify repository has been attached to the
// Gin context by SessionMiddleware. Aborts with 401 Unauthorized if the repository
// is missing or nil.
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
