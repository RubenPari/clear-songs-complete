// Package http configures the HTTP transport layer. It wires up Gin routes,
// middleware, and controllers using the DI container.
package http

import (
	"github.com/RubenPari/clear-songs/internal/infrastructure/di"
	"github.com/RubenPari/clear-songs/internal/infrastructure/transport/http/handlers"
	"github.com/RubenPari/clear-songs/internal/infrastructure/transport/http/middleware"
	"github.com/gin-gonic/gin"
)

// SetUpRoutes configures all HTTP routes on the Gin engine using the DI container
// for dependency injection. It registers health, auth, track, and playlist endpoints
// with appropriate middleware chains.
func SetUpRoutes(server *gin.Engine, container *di.Container) {
	setUpHealthRoute(server)

	server.Use(middleware.SessionMiddleware(
		container.SpotifyRepo,
		container.CacheRepo,
	))
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

// setUpHealthRoute registers the unauthenticated health check endpoint.
func setUpHealthRoute(server *gin.Engine) {
	server.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}
