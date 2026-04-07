package track

import (
	"context"
	spotifyAPI "github.com/zmb3/spotify"
)

// TrackRepository defines the interface for track-related operations
type TrackRepository interface {
	// GetTrackSummary retrieves a summary of tracks grouped by artist
	GetTrackSummary(ctx context.Context, min, max int) ([]ArtistSummary, error)
	
	// GetTracksByArtist retrieves all tracks from a specific artist
	GetTracksByArtist(ctx context.Context, artistID spotifyAPI.ID) ([]spotifyAPI.SavedTrack, error)
	
	// DeleteTracksByArtist removes all tracks from a specific artist
	DeleteTracksByArtist(ctx context.Context, artistID spotifyAPI.ID) error
	
	// DeleteTracksByRange removes tracks within a count range
	DeleteTracksByRange(ctx context.Context, min, max int) error
}
