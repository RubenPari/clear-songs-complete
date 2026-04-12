package track

import (
	"context"

	spotifyAPI "github.com/zmb3/spotify"
)

// getUserTracks retrieves tracks from cache or API.
func (uc *GetTrackSummaryUseCase) getUserTracks(ctx context.Context) ([]spotifyAPI.SavedTrack, error) {
	return getUserTracks(ctx, uc.spotifyRepo, uc.cacheRepo)
}
