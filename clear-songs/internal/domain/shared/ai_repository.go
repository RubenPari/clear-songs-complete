package shared

import "context"

// AIArtistLookup identifies an artist for batch genre resolution. Key must match
// the stable map key used in track summary (Spotify artist ID or synthetic key when ID is missing).
type AIArtistLookup struct {
	Key  string
	Name string
}

// AIRepository defines the interface for AI-powered operations.
type AIRepository interface {
	// ResolveArtistGenres returns raw genre strings from the model keyed by Lookup.Key
	// (lowercase trimmed). The application layer maps them with NormalizeAIGenreLabel and ResolveGenre.
	ResolveArtistGenres(ctx context.Context, lookups []AIArtistLookup) (map[string]string, error)
}
