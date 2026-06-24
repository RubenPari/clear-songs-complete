package handlers

import (
	"github.com/RubenPari/clear-songs/internal/application/track"
	"github.com/gin-gonic/gin"
	spotifyAPI "github.com/zmb3/spotify"
)


// TrackController is the refactored track controller
type TrackController struct {
	BaseController
	getTrackSummaryUseCase     *track.GetTrackSummaryUseCase
	deleteTracksByArtistUseCase *track.DeleteTracksByArtistUseCase
	deleteTracksByRangeUseCase  *track.DeleteTracksByRangeUseCase
	deleteTrackUseCase          *track.DeleteTrackUseCase
	getTracksByArtistUseCase    *track.GetTracksByArtistUseCase
	invalidateLibraryCacheUseCase *track.InvalidateLibraryCacheUseCase
}

// Creates track controller.
func NewTrackController(
	getTrackSummaryUseCase *track.GetTrackSummaryUseCase,
	deleteTracksByArtistUseCase *track.DeleteTracksByArtistUseCase,
	deleteTracksByRangeUseCase *track.DeleteTracksByRangeUseCase,
	getTracksByArtistUseCase *track.GetTracksByArtistUseCase,
	deleteTrackUseCase *track.DeleteTrackUseCase,
	invalidateLibraryCacheUseCase *track.InvalidateLibraryCacheUseCase,
) *TrackController {
	return &TrackController{
		getTrackSummaryUseCase:      getTrackSummaryUseCase,
		deleteTracksByArtistUseCase: deleteTracksByArtistUseCase,
		deleteTracksByRangeUseCase:  deleteTracksByRangeUseCase,
		deleteTrackUseCase:          deleteTrackUseCase,
		getTracksByArtistUseCase:    getTracksByArtistUseCase,
		invalidateLibraryCacheUseCase: invalidateLibraryCacheUseCase,
	}
}

// Invalidates library cache.
func (tc *TrackController) InvalidateLibraryCache(c *gin.Context) {
	ctx := c.Request.Context()
	msg, err := tc.invalidateLibraryCacheUseCase.Execute(ctx)
	if err != nil {
		tc.JSONInternalError(c, "Failed to invalidate library cache")
		return
	}
	tc.JSONSuccess(c, gin.H{"message": msg})
}

// Fetches track summary.
func (tc *TrackController) GetTrackSummary(c *gin.Context) {
	var req track.RangeRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		tc.JSONValidationError(c, "Invalid min or max parameters")
		return
	}

	min, max, errMsg := track.ValidateRangeQuery(&req)
	if errMsg != "" {
		tc.JSONValidationError(c, errMsg)
		return
	}

	ctx := c.Request.Context()
	result, err := tc.getTrackSummaryUseCase.Execute(ctx, min, max, req.Genre)
	if err != nil {
		tc.HandleDomainError(c, err)
		return
	}

	response := make([]track.ArtistSummaryResponse, len(result))
	for i, artist := range result {
		response[i] = track.NewArtistSummaryResponse(artist)
	}

	tc.JSONSuccess(c, response)
}

// Fetches tracks by artist.
func (tc *TrackController) GetTracksByArtist(c *gin.Context) {
	// Get artist ID from URL
	idArtistString := c.Param("id_artist")
	if idArtistString == "" {
		tc.JSONValidationError(c, "Artist ID is required")
		return
	}

	artistID := spotifyAPI.ID(idArtistString)

	// Execute use case
	ctx := c.Request.Context()
	tracks, err := tc.getTracksByArtistUseCase.Execute(ctx, artistID)
	if err != nil {
		tc.HandleDomainError(c, err)
		return
	}

	response := make([]track.TrackResponse, len(tracks))
	for i, t := range tracks {
		response[i] = track.NewTrackResponse(t)
	}

	tc.JSONSuccess(c, response)
}

// Deletes track by artist.
func (tc *TrackController) DeleteTrackByArtist(c *gin.Context) {
	// Get artist ID from URL
	idArtistString := c.Param("id_artist")
	if idArtistString == "" {
		tc.JSONValidationError(c, "Artist ID is required")
		return
	}

	artistID := spotifyAPI.ID(idArtistString)

	// Execute use case
	ctx := c.Request.Context()
	if err := tc.deleteTracksByArtistUseCase.Execute(ctx, artistID); err != nil {
		tc.HandleDomainError(c, err)
		return
	}

	tc.JSONSuccess(c, gin.H{"message": "Tracks deleted successfully"})
}

// Deletes track.
func (tc *TrackController) DeleteTrack(c *gin.Context) {
	// Get track ID from URL
	idTrackString := c.Param("id_track")
	if idTrackString == "" {
		tc.JSONValidationError(c, "Track ID is required")
		return
	}

	trackID := spotifyAPI.ID(idTrackString)

	// Execute use case
	ctx := c.Request.Context()
	if err := tc.deleteTrackUseCase.Execute(ctx, trackID); err != nil {
		tc.HandleDomainError(c, err)
		return
	}

	tc.JSONSuccess(c, gin.H{"message": "Track deleted successfully"})
}

// Deletes track by range.
func (tc *TrackController) DeleteTrackByRange(c *gin.Context) {
	var req track.RangeRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		tc.JSONValidationError(c, "Invalid min or max parameters")
		return
	}

	// At least one parameter must be provided for a destructive action
	if c.Query("min") == "" && c.Query("max") == "" {
		tc.JSONValidationError(c, "At least one of min or max must be provided")
		return
	}

	min, max, errMsg := track.ValidateRangeQuery(&req)
	if errMsg != "" {
		tc.JSONValidationError(c, errMsg)
		return
	}

	// Execute use case
	ctx := c.Request.Context()
	if err := tc.deleteTracksByRangeUseCase.Execute(ctx, min, max); err != nil {
		tc.HandleDomainError(c, err)
		return
	}

	tc.JSONSuccess(c, gin.H{"message": "Tracks deleted successfully"})
}
