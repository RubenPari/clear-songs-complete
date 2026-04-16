package track

import domainTrack "github.com/RubenPari/clear-songs/internal/domain/track"

// Filter summary by range.
func FilterSummaryByRange(tracks []domainTrack.ArtistSummary, min, max int) []domainTrack.ArtistSummary {
	filtered := make([]domainTrack.ArtistSummary, 0, len(tracks))

	for _, artist := range tracks {
		if min > 0 && artist.Count < min {
			continue
		}
		if max > 0 && artist.Count > max {
			continue
		}

		filtered = append(filtered, artist)
	}

	return filtered
}
