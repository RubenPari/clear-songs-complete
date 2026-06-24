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

// PlaylistController is the refactored playlist controller using dependency injection
type PlaylistController struct {
	BaseController
	getUserPlaylistsUseCase         *playlist.GetUserPlaylistsUseCase
	deletePlaylistTracksUseCase     *playlist.DeletePlaylistTracksUseCase
	deletePlaylistAndLibraryUseCase *playlist.DeletePlaylistAndLibraryTracksUseCase
}

// Creates playlist controller.
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

// Fetches user playlists.
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

// Deletes all playlist tracks.
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

// Deletes all playlist and user tracks.
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
