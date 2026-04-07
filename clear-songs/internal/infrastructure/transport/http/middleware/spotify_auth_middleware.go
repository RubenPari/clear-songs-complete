package middleware

import (
	"net/http"

	"github.com/RubenPari/clear-songs/internal/application/shared/dto"
	"github.com/RubenPari/clear-songs/internal/domain/shared"
	"github.com/gin-gonic/gin"
	spotifyAPI "github.com/zmb3/spotify"
)

// SpotifyAuthMiddleware requires Spotify authentication
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

		// Get client from repository (for backward compatibility)
		if spotifyRepoImpl, ok := spotifyRepo.(interface{ GetClient() *spotifyAPI.Client }); ok {
			client := spotifyRepoImpl.GetClient()
			if client == nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, dto.UnauthorizedErr())
				return
			}
			c.Set("spotifyClient", client)
		}

		// Also store the repository for direct use
		c.Set("spotifyRepository", spotifyRepo)
		c.Next()
	}
}
