package track

import spotifyAPI "github.com/zmb3/spotify"

// NewTrackResponse maps a Spotify saved track to the API response DTO.
func NewTrackResponse(t spotifyAPI.SavedTrack) TrackResponse {
	artists := make([]string, len(t.Artists))
	for i, artist := range t.Artists {
		artists[i] = artist.Name
	}

	imageURL := ""
	if len(t.Album.Images) > 0 {
		imageURL = getBestImageURL(t.Album.Images, 300)
	}

	spotifyURL := ""
	if url, exists := t.ExternalURLs["spotify"]; exists {
		spotifyURL = url
	}

	return TrackResponse{
		ID:         t.ID.String(),
		Name:       t.Name,
		Artists:    artists,
		Album:      t.Album.Name,
		Duration:   t.Duration,
		ImageURL:   imageURL,
		SpotifyURL: spotifyURL,
	}
}

// getBestImageURL returns the smallest image URL that fits within maxWidth.
// It iterates from the smallest image upward, which is a common ordering.
func getBestImageURL(images []spotifyAPI.Image, maxWidth int) string {
	if len(images) == 0 {
		return ""
	}

	for i := len(images) - 1; i >= 0; i-- {
		if images[i].Width <= maxWidth || i == 0 {
			return images[i].URL
		}
	}

	return ""
}
