package shared

import "context"

// AIRepository defines the interface for AI-powered operations
type AIRepository interface {
	// ResolveArtistGenre asks the AI model to determine the genre for an artist
	ResolveArtistGenre(ctx context.Context, artistName string) (string, error)
}
