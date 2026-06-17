package playlist

import (
	spotifyhelpers "github.com/RubenPari/clear-songs/internal/infrastructure/external/spotify/helpers"
	spotifyAPI "github.com/zmb3/spotify"
)

// NewPlaylistResponse maps a Spotify simple playlist to the API response DTO.
func NewPlaylistResponse(p spotifyAPI.SimplePlaylist) PlaylistResponse {
	return PlaylistResponse{
		ID:       p.ID.String(),
		Name:     p.Name,
		ImageURL: spotifyhelpers.GetMediumImage(p.Images),
	}
}
