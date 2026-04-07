package shared

import (
	"context"

	spotifyAPI "github.com/zmb3/spotify"
)

// SpotifyRepository defines the interface for Spotify API operations
type SpotifyRepository interface {
	// GetCurrentUser retrieves the current authenticated user
	GetCurrentUser(ctx context.Context) (*spotifyAPI.PrivateUser, error)

	// GetUserTracks retrieves all tracks saved by the user
	GetUserTracks(ctx context.Context, limit, offset int) ([]spotifyAPI.SavedTrack, error)

	// GetAllUserTracks retrieves all user tracks with pagination
	GetAllUserTracks(ctx context.Context) ([]spotifyAPI.SavedTrack, error)

	// GetTracksByArtist filters tracks by artist ID
	GetTracksByArtist(ctx context.Context, artistID spotifyAPI.ID, tracks []spotifyAPI.SavedTrack) ([]spotifyAPI.SavedTrack, error)

	// GetTrackIDsByArtist returns only track IDs for an artist
	GetTrackIDsByArtist(ctx context.Context, artistID spotifyAPI.ID, tracks []spotifyAPI.SavedTrack) ([]spotifyAPI.ID, error)

	// DeleteTracksFromLibrary removes tracks from user's library
	DeleteTracksFromLibrary(ctx context.Context, trackIDs []spotifyAPI.ID) error

	// GetPlaylist retrieves a playlist by ID
	GetPlaylist(ctx context.Context, playlistID spotifyAPI.ID) (*spotifyAPI.FullPlaylist, error)

	// GetPlaylistTracks retrieves all tracks from a playlist
	GetPlaylistTracks(ctx context.Context, playlistID spotifyAPI.ID, limit, offset int) ([]spotifyAPI.PlaylistTrack, error)

	// GetAllPlaylistTracks retrieves all tracks from a playlist with pagination
	GetAllPlaylistTracks(ctx context.Context, playlistID spotifyAPI.ID) ([]spotifyAPI.PlaylistTrack, error)

	// DeletePlaylistTracks removes tracks from a playlist
	DeletePlaylistTracks(ctx context.Context, playlistID spotifyAPI.ID, trackIDs []spotifyAPI.ID) error

	// GetUserPlaylists retrieves all playlists owned or followed by the user
	GetUserPlaylists(ctx context.Context, limit, offset int) ([]spotifyAPI.SimplePlaylist, error)

	// GetAllUserPlaylists retrieves all user playlists with pagination
	GetAllUserPlaylists(ctx context.Context) ([]spotifyAPI.SimplePlaylist, error)

	// GetArtist retrieves artist information
	GetArtist(ctx context.Context, artistID spotifyAPI.ID) (*spotifyAPI.FullArtist, error)

	// GetArtists retrieves multiple artists in batch (up to 50 per call)
	GetArtists(ctx context.Context, artistIDs []spotifyAPI.ID) ([]*spotifyAPI.FullArtist, error)

	// GetTrack retrieves track information
	GetTrack(ctx context.Context, trackID spotifyAPI.ID) (*spotifyAPI.FullTrack, error)

	// SetAccessToken sets the OAuth token for authenticated requests
	SetAccessToken(token interface{}) error
}
