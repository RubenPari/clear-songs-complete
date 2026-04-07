package handlers

import (
	"github.com/RubenPari/clear-songs/internal/application/playlist"
	"github.com/RubenPari/clear-songs/internal/domain/shared/utils"
	"github.com/gin-gonic/gin"
	spotifyAPI "github.com/zmb3/spotify"
)

// PlaylistRequest validates the incoming query parameters
type PlaylistRequest struct {
	ID string `form:"id" binding:"required"`
}

// PlaylistControllerRefactored is the refactored playlist controller using dependency injection
type PlaylistControllerRefactored struct {
	BaseController
	getUserPlaylistsUC         *playlist.GetUserPlaylistsUseCase
	deletePlaylistTracksUC     *playlist.DeletePlaylistTracksUseCase
	deletePlaylistAndLibraryUC *playlist.DeletePlaylistAndLibraryTracksUseCase
}

// NewPlaylistController creates a new PlaylistController
func NewPlaylistController(
	getUserPlaylistsUC *playlist.GetUserPlaylistsUseCase,
	deletePlaylistTracksUC *playlist.DeletePlaylistTracksUseCase,
	deletePlaylistAndLibraryUC *playlist.DeletePlaylistAndLibraryTracksUseCase,
) *PlaylistControllerRefactored {
	return &PlaylistControllerRefactored{
		getUserPlaylistsUC:         getUserPlaylistsUC,
		deletePlaylistTracksUC:     deletePlaylistTracksUC,
		deletePlaylistAndLibraryUC: deletePlaylistAndLibraryUC,
	}
}

// GetUserPlaylists handles GET /playlist/list
func (pc *PlaylistControllerRefactored) GetUserPlaylists(c *gin.Context) {
	ctx := c.Request.Context()
	playlists, err := pc.getUserPlaylistsUC.Execute(ctx)
	if err != nil {
		pc.HandleDomainError(c, err)
		return
	}

	// Convert to response format
	var response []playlist.PlaylistResponse
	for _, p := range playlists {
		imageURL := utils.GetMediumImage(p.Images)

		response = append(response, playlist.PlaylistResponse{
			ID:       p.ID.String(),
			Name:     p.Name,
			ImageURL: imageURL,
		})
	}

	pc.JSONSuccess(c, response)
}

// DeleteAllPlaylistTracks handles DELETE /playlist/delete-tracks
func (pc *PlaylistControllerRefactored) DeleteAllPlaylistTracks(c *gin.Context) {
	var req PlaylistRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		pc.JSONValidationError(c, "Playlist id is required")
		return
	}

	playlistID := spotifyAPI.ID(req.ID)
	ctx := c.Request.Context()

	if err := pc.deletePlaylistTracksUC.Execute(ctx, playlistID); err != nil {
		pc.HandleDomainError(c, err)
		return
	}

	pc.JSONSuccess(c, gin.H{"message": "Tracks deleted successfully"})
}

// DeleteAllPlaylistAndUserTracks handles DELETE /playlist/delete-tracks-and-library
func (pc *PlaylistControllerRefactored) DeleteAllPlaylistAndUserTracks(c *gin.Context) {
	var req PlaylistRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		pc.JSONValidationError(c, "Playlist id is required")
		return
	}

	playlistID := spotifyAPI.ID(req.ID)
	ctx := c.Request.Context()

	if err := pc.deletePlaylistAndLibraryUC.Execute(ctx, playlistID); err != nil {
		pc.HandleDomainError(c, err)
		return
	}

	pc.JSONSuccess(c, gin.H{"message": "Tracks deleted successfully"})
}
