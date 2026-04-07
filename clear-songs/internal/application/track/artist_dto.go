package track

// ArtistSummary represents an artist summary in API responses
type ArtistSummary struct {
	Id       string   `json:"id"`
	Name     string   `json:"name"`
	Count    int      `json:"count"`
	ImageURL string   `json:"image_url,omitempty"`
	Genres   []string `json:"genres,omitempty"`
	Genre    string   `json:"genre,omitempty"`
}
