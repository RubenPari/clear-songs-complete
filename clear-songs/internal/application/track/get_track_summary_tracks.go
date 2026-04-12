package track

import (
	"context"
	"time"

	spotifyAPI "github.com/zmb3/spotify"
)

// getUserTracks retrieves tracks from cache or API.
func (uc *GetTrackSummaryUseCase) getUserTracks(ctx context.Context) ([]spotifyAPI.SavedTrack, error) {
	if uc.cacheRepo != nil {
		cached, err := uc.cacheRepo.GetUserTracks(ctx)
		if err == nil && len(cached) > 0 {
			return cached, nil
		}
	}

	tracks, err := uc.spotifyRepo.GetAllUserTracks(ctx)
	if err != nil {
		return nil, err
	}

	if uc.cacheRepo != nil {
		_ = uc.cacheRepo.SetUserTracks(ctx, tracks, 5*time.Minute)
	}

	return tracks, nil
}
