package track

import (
	"context"
	"sort"
)

// GetArtistGenresDebug returns every primary artist in the user library with Spotify genre tags.
// No min/max/genre filter; intended for debugging and tuning genre mapping.
func (uc *GetTrackSummaryUseCase) GetArtistGenresDebug(ctx context.Context) ([]ArtistGenresDebugEntry, error) {
	tracks, err := uc.getUserTracks(ctx)
	if err != nil {
		return nil, err
	}

	artistMap := groupTracksByPrimaryArtist(tracks)
	artistIDs := collectArtistIDs(artistMap)
	details := uc.fetchArtistDetails(ctx, artistIDs)

	out := make([]ArtistGenresDebugEntry, 0, len(artistMap))
	for _, data := range artistMap {
		_, genres := extractArtistMetadata(data.id, details)
		if genres == nil {
			genres = []string{}
		}
		out = append(out, ArtistGenresDebugEntry{
			ID:         data.id,
			Name:       data.name,
			TrackCount: data.count,
			Genres:     genres,
		})
	}

	sort.Slice(out, func(i, j int) bool {
		if out[i].Name != out[j].Name {
			return out[i].Name < out[j].Name
		}
		return out[i].ID < out[j].ID
	})

	return out, nil
}
