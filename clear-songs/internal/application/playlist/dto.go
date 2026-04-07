package playlist

// PlaylistResponse represents a playlist in API responses
type PlaylistResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	ImageURL string `json:"image_url,omitempty"`
}
