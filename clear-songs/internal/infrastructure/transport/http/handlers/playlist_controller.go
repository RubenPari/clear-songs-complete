package handlers

import (
	"github.com/RubenPari/clear-songs/internal/application/playlist"
	"github.com/gin-gonic/gin"
	spotifyAPI "github.com/zmb3/spotify"
)


// PlaylistRequest validates the incoming query parameters
type PlaylistRequest struct {
	ID string `form:"id" binding:"required"`
}

// PlaylistController handles playlist-related HTTP endpoints.
type PlaylistController struct {
	BaseController
	getUserPlaylistsUseCase         *playlist.GetUserPlaylistsUseCase
	deletePlaylistTracksUseCase     *playlist.DeletePlaylistTracksUseCase
	deletePlaylistAndLibraryUseCase *playlist.DeletePlaylistAndLibraryTracksUseCase
}

// NewPlaylistController creates a playlist controller with the given use cases.
func NewPlaylistController(
	getUserPlaylistsUseCase *playlist.GetUserPlaylistsUseCase,
	deletePlaylistTracksUseCase *playlist.DeletePlaylistTracksUseCase,
	deletePlaylistAndLibraryUseCase *playlist.DeletePlaylistAndLibraryTracksUseCase,
) *PlaylistController {
	return &PlaylistController{
		getUserPlaylistsUseCase:         getUserPlaylistsUseCase,
		deletePlaylistTracksUseCase:     deletePlaylistTracksUseCase,
		deletePlaylistAndLibraryUseCase: deletePlaylistAndLibraryUseCase,
	}
}

// GetUserPlaylists retrieves all playlists owned by the authenticated user.
func (pc *PlaylistController) GetUserPlaylists(c *gin.Context) {
	ctx := c.Request.Context()
	playlists, err := pc.getUserPlaylistsUseCase.Execute(ctx)
	if err != nil {
		pc.HandleDomainError(c, err)
		return
	}

	response := make([]playlist.PlaylistResponse, len(playlists))
	for i, p := range playlists {
		response[i] = playlist.NewPlaylistResponse(p)
	}

	pc.JSONSuccess(c, response)
}

// DeleteAllPlaylistTracks removes all tracks from the specified playlist.
// Requires the playlist ID as a query parameter.
func (pc *PlaylistController) DeleteAllPlaylistTracks(c *gin.Context) {
	var req PlaylistRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		pc.JSONValidationError(c, "Playlist id is required")
		return
	}

	playlistID := spotifyAPI.ID(req.ID)
	ctx := c.Request.Context()

	if err := pc.deletePlaylistTracksUseCase.Execute(ctx, playlistID); err != nil {
		pc.HandleDomainError(c, err)
		return
	}

	pc.JSONSuccess(c, gin.H{"message": "Tracks deleted successfully"})
}

// DeleteAllPlaylistAndUserTracks removes all tracks from the specified playlist
// and also removes them from the user's saved library.
func (pc *PlaylistController) DeleteAllPlaylistAndUserTracks(c *gin.Context) {
	var req PlaylistRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		pc.JSONValidationError(c, "Playlist id is required")
		return
	}

	playlistID := spotifyAPI.ID(req.ID)
	ctx := c.Request.Context()

	if err := pc.deletePlaylistAndLibraryUseCase.Execute(ctx, playlistID); err != nil {
		pc.HandleDomainError(c, err)
		return
	}

	pc.JSONSuccess(c, gin.H{"message": "Tracks deleted successfully"})
}
