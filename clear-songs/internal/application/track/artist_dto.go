package track

import domainTrack "github.com/RubenPari/clear-songs/internal/domain/track"

// ArtistSummaryResponse represents an artist summary in API responses.
type ArtistSummaryResponse struct {
	Id       string   `json:"id"`
	Name     string   `json:"name"`
	Count    int      `json:"count"`
	ImageURL string   `json:"image_url,omitempty"`
	Genres   []string `json:"genres,omitempty"`
	Genre    string   `json:"genre,omitempty"`
}

// NewArtistSummaryResponse maps a domain ArtistSummary to its API representation.
func NewArtistSummaryResponse(artist domainTrack.ArtistSummary) ArtistSummaryResponse {
	return ArtistSummaryResponse{
		Id:       artist.ID,
		Name:     artist.Name,
		Count:    artist.Count,
		ImageURL: artist.ImageURL,
		Genres:   artist.Genres,
		Genre:    artist.Genre,
	}
}
