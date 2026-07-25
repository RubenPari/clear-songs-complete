package track

import (
	"context"

	spotifyAPI "github.com/zmb3/spotify"
)

// getUserTracks delegates to the shared helper that reads from cache or fetches
// from the Spotify API. This method is specific to GetTrackSummaryUseCase.
func (uc *GetTrackSummaryUseCase) getUserTracks(ctx context.Context) ([]spotifyAPI.SavedTrack, error) {
	return getUserTracks(ctx, uc.spotifyRepo, uc.cacheRepo)
}
