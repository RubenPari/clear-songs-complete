package playlist

import (
	"context"
	spotifyAPI "github.com/zmb3/spotify"
)

// PlaylistRepository defines the interface for playlist-related operations
type PlaylistRepository interface {
	// GetUserPlaylists retrieves all playlists owned or followed by the user
	GetUserPlaylists(ctx context.Context) ([]spotifyAPI.SimplePlaylist, error)
	
	// DeletePlaylistTracks removes all tracks from a playlist
	DeletePlaylistTracks(ctx context.Context, playlistID spotifyAPI.ID) error
	
	// DeletePlaylistAndLibraryTracks removes tracks from both playlist and user library
	DeletePlaylistAndLibraryTracks(ctx context.Context, playlistID spotifyAPI.ID) error
}
