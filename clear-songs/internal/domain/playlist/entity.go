package playlist

import "github.com/RubenPari/clear-songs/internal/domain/track"

// Playlist represents a playlist entity
type Playlist struct {
	ID       string
	Name     string
	ImageURL string
	Tracks   []track.Track
}
