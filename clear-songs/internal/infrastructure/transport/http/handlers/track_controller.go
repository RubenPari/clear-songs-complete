package handlers

import (
	"github.com/RubenPari/clear-songs/internal/application/track"
	"github.com/gin-gonic/gin"
	spotifyAPI "github.com/zmb3/spotify"
)


// TrackController handles track-related HTTP endpoints.
type TrackController struct {
	BaseController
	getTrackSummaryUseCase     *track.GetTrackSummaryUseCase
	deleteTracksByArtistUseCase *track.DeleteTracksByArtistUseCase
	deleteTracksByRangeUseCase  *track.DeleteTracksByRangeUseCase
	deleteTrackUseCase          *track.DeleteTrackUseCase
	getTracksByArtistUseCase    *track.GetTracksByArtistUseCase
	invalidateLibraryCacheUseCase *track.InvalidateLibraryCacheUseCase
}

// NewTrackController creates a track controller with the given use cases.
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

// InvalidateLibraryCache clears the cached user tracks and summary entries.
func (tc *TrackController) InvalidateLibraryCache(c *gin.Context) {
	ctx := c.Request.Context()
	msg, err := tc.invalidateLibraryCacheUseCase.Execute(ctx)
	if err != nil {
		tc.JSONInternalError(c, "Failed to invalidate library cache")
		return
	}
	tc.JSONSuccess(c, gin.H{"message": msg})
}

// GetTrackSummary retrieves a summary of the user's saved tracks grouped by artist.
// Supports optional min, max, and genre query parameters for filtering.
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

// GetTracksByArtist retrieves all saved tracks by the specified artist.
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

// DeleteTrackByArtist removes all saved tracks by the specified artist from the library.
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

// DeleteTrack removes a single track from the user's library.
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

// DeleteTrackByRange removes all saved tracks whose primary artist has a track count
// within the specified range. At least one of min or max must be provided.
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
