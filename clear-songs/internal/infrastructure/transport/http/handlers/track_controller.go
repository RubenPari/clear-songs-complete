package handlers

import (
	"github.com/RubenPari/clear-songs/internal/application/track"
	"github.com/RubenPari/clear-songs/internal/domain/shared"
	"github.com/RubenPari/clear-songs/internal/domain/shared/utils"
	"github.com/gin-gonic/gin"
	spotifyAPI "github.com/zmb3/spotify"
)

// TrackController is the refactored track controller
type TrackController struct {
	BaseController
	cacheRepo              shared.CacheRepository
	getTrackSummaryUseCase *track.GetTrackSummaryUseCase
	deleteTracksByArtistUC *track.DeleteTracksByArtistUseCase
	deleteTracksByRangeUC  *track.DeleteTracksByRangeUseCase
	deleteTrackUC          *track.DeleteTrackUseCase
	getTracksByArtistUC    *track.GetTracksByArtistUseCase
}

// NewTrackController creates a new TrackController
func NewTrackController(
	cacheRepo shared.CacheRepository,
	getTrackSummaryUseCase *track.GetTrackSummaryUseCase,
	deleteTracksByArtistUC *track.DeleteTracksByArtistUseCase,
	deleteTracksByRangeUC *track.DeleteTracksByRangeUseCase,
	getTracksByArtistUC *track.GetTracksByArtistUseCase,
	deleteTrackUC *track.DeleteTrackUseCase,
) *TrackController {
	return &TrackController{
		cacheRepo:              cacheRepo,
		getTrackSummaryUseCase: getTrackSummaryUseCase,
		deleteTracksByArtistUC: deleteTracksByArtistUC,
		deleteTracksByRangeUC:  deleteTracksByRangeUC,
		deleteTrackUC:          deleteTrackUC,
		getTracksByArtistUC:    getTracksByArtistUC,
	}
}

// InvalidateLibraryCache handles POST /track/library-cache/invalidate — clears user tracks
// and derived track-summary keys in Redis so the next GET /track/summary is recomputed.
func (tc *TrackController) InvalidateLibraryCache(c *gin.Context) {
	if tc.cacheRepo == nil {
		tc.JSONSuccess(c, gin.H{"message": "No cache configured"})
		return
	}
	ctx := c.Request.Context()
	if err := tc.cacheRepo.InvalidateUserTracks(ctx); err != nil {
		tc.JSONInternalError(c, "Failed to invalidate library cache")
		return
	}
	tc.JSONSuccess(c, gin.H{"message": "Library cache invalidated"})
}

// GetTrackSummary handles GET /track/summary
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

	var response []track.ArtistSummary
	for _, artist := range result {
		response = append(response, track.ArtistSummary{
			Id:       artist.ID,
			Name:     artist.Name,
			Count:    artist.Count,
			ImageURL: artist.ImageURL,
			Genres:   artist.Genres,
			Genre:    artist.Genre,
		})
	}

	tc.JSONSuccess(c, response)
}

// GetTracksByArtist handles GET /track/by-artist/:id_artist
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
	tracks, err := tc.getTracksByArtistUC.Execute(ctx, artistID)
	if err != nil {
		tc.HandleDomainError(c, err)
		return
	}

	// Convert to response format
	var response []track.TrackResponse
	for _, t := range tracks {
		artists := make([]string, len(t.Artists))
		for i, artist := range t.Artists {
			artists[i] = artist.Name
		}

		imageURL := utils.GetMediumImage(t.Album.Images)

		spotifyURL := ""
		if url, exists := t.ExternalURLs["spotify"]; exists {
			spotifyURL = url
		}

		response = append(response, track.TrackResponse{
			ID:         t.ID.String(),
			Name:       t.Name,
			Artists:    artists,
			Album:      t.Album.Name,
			Duration:   t.Duration,
			ImageURL:   imageURL,
			SpotifyURL: spotifyURL,
		})
	}

	tc.JSONSuccess(c, response)
}

// DeleteTrackByArtist handles DELETE /track/by-artist/:id_artist
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
	if err := tc.deleteTracksByArtistUC.Execute(ctx, artistID); err != nil {
		tc.HandleDomainError(c, err)
		return
	}

	tc.JSONSuccess(c, gin.H{"message": "Tracks deleted successfully"})
}

// DeleteTrack handles DELETE /track/:id_track
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
	if err := tc.deleteTrackUC.Execute(ctx, trackID); err != nil {
		tc.HandleDomainError(c, err)
		return
	}

	tc.JSONSuccess(c, gin.H{"message": "Track deleted successfully"})
}

// DeleteTrackByRange handles DELETE /track/by-range
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
	if err := tc.deleteTracksByRangeUC.Execute(ctx, min, max); err != nil {
		tc.HandleDomainError(c, err)
		return
	}

	tc.JSONSuccess(c, gin.H{"message": "Tracks deleted successfully"})
}
