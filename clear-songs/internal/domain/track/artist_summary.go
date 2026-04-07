package track

// ArtistSummary represents a summary of tracks by artist
type ArtistSummary struct {
	ID       string
	Name     string
	Count    int
	ImageURL string
	Genres   []string
	Genre    string
}
