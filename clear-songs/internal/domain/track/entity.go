package track

// Track represents a track entity in the domain
type Track struct {
	ID       string
	Name     string
	Artists  []Artist
	Album    Album
	Duration int
	ImageURL string
	SpotifyURL string
}

// Artist represents an artist entity
type Artist struct {
	ID       string
	Name     string
	ImageURL string
}

// Album represents an album entity
type Album struct {
	ID       string
	Name     string
	ImageURL string
}
