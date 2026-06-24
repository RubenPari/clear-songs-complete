package http

import (
	"github.com/RubenPari/clear-songs/internal/infrastructure/di"
	"github.com/RubenPari/clear-songs/internal/infrastructure/transport/http/handlers"
	"github.com/RubenPari/clear-songs/internal/infrastructure/transport/http/middleware"
	"github.com/gin-gonic/gin"
)

/**
 * SetUpRoutes configures all HTTP routes using dependency injection
 *
 * This version uses the DI container to inject dependencies into controllers
 * and middleware, eliminating the need for global variables.
 *
 * @param server - The Gin engine instance to configure routes on
 * @param container - The dependency injection container
 */
// Registers HTTP routes.
func SetUpRoutes(server *gin.Engine, container *di.Container) {
	server.Use(middleware.SessionMiddleware(
		container.SpotifyRepo,
		container.CacheRepo,
	))
	server.Use(middleware.CacheInvalidationMiddleware())

	server.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"status":  "error",
			"message": "not found path",
		})
	})

	trackController := handlers.NewTrackController(
		container.GetTrackSummaryUseCase,
		container.DeleteTracksByArtistUseCase,
		container.DeleteTracksByRangeUseCase,
		container.GetTracksByArtistUseCase,
		container.DeleteTrackUseCase,
		container.InvalidateLibraryCacheUseCase,
	)

	track := server.Group("/track")
	{
		track.POST("/library-cache/invalidate",
			middleware.SpotifyAuthMiddleware(),
			trackController.InvalidateLibraryCache)
		track.GET("/summary",
			middleware.SpotifyAuthMiddleware(),
			trackController.GetTrackSummary)
		track.GET("/by-artist/:id_artist",
			middleware.SpotifyAuthMiddleware(),
			trackController.GetTracksByArtist)
		track.DELETE("/by-artist/:id_artist",
			middleware.SpotifyAuthMiddleware(),
			trackController.DeleteTrackByArtist)
		track.DELETE("/:id_track",
			middleware.SpotifyAuthMiddleware(),
			trackController.DeleteTrack)
		track.DELETE("/by-range",
			middleware.SpotifyAuthMiddleware(),
			trackController.DeleteTrackByRange)
	}

	authController := handlers.NewAuthController(
		container.LoginUseCase,
		container.CallbackUseCase,
		container.LogoutUseCase,
		container.IsAuthUseCase,
	)

	auth := server.Group("/auth")
	{
		auth.GET("/login", authController.Login)
		auth.GET("/callback", authController.Callback)
		auth.GET("/logout", authController.Logout)
		auth.GET("/is-auth", authController.IsAuth)
	}

	playlistController := handlers.NewPlaylistController(
		container.GetUserPlaylistsUseCase,
		container.DeletePlaylistTracksUseCase,
		container.DeletePlaylistAndLibraryUseCase,
	)

	playlist := server.Group("/playlist")
	{
		playlist.GET("/list",
			middleware.SpotifyAuthMiddleware(),
			playlistController.GetUserPlaylists)
		playlist.DELETE("/delete-tracks",
			middleware.SpotifyAuthMiddleware(),
			playlistController.DeleteAllPlaylistTracks)
		playlist.DELETE("/delete-tracks-and-library",
			middleware.SpotifyAuthMiddleware(),
			playlistController.DeleteAllPlaylistAndUserTracks)
	}
}
