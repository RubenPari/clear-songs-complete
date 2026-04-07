package shared

import (
	spotifyAPI "github.com/zmb3/spotify"
)

// DatabaseRepository defines the interface for database operations
type DatabaseRepository interface {
	// SaveTracksBackup saves tracks to database as backup
	SaveTracksBackup(tracks []spotifyAPI.PlaylistTrack) error

	// SaveFullTracksBackup saves full tracks to database as backup
	SaveFullTracksBackup(tracks []spotifyAPI.FullTrack) error
}
