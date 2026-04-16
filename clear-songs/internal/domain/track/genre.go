package track

import "strings"

// genreMapping defines canonical genres and their keyword matchers.
// Order matters: first match wins. Keywords are checked with strings.Contains
// on the lowercased Spotify genre, longest keywords first within each entry.
var genreMapping = []struct {
	Canonical string
	Keywords  []string
}{
	// Folk must come before Rock (folk rock contains "rock")
	{
		Canonical: "Folk",
		Keywords: []string{
			"folk rock", "folk", "singer-songwriter", "roots rock",
		},
	},
	// Metal must come before Hip Hop (rap metal contains "rap")
	{
		Canonical: "Metal",
		Keywords: []string{
			"alternative metal", "nu metal", "thrash metal", "heavy metal",
			"glam metal", "doom metal", "stoner rock", "rap metal", "metal",
		},
	},
	{
		Canonical: "Hip Hop",
		Keywords: []string{
			"hip hop", "hip-hop", "rap", "trap", "boom bap", "gangster",
			"grime", "drill", "horrorcore", "crunk", "bounce",
			"christian hip hop", "country hip hop", "boogie",
		},
	},
	{
		Canonical: "Electronic",
		Keywords: []string{
			"progressive house", "tropical house", "electro house",
			"slap house", "big room", "drum and bass", "lo-fi hip hop",
			"lo-fi", "edm", "electronic", "house", "trance", "dubstep",
			"techno", "dnb", "moombahton", "electronica", "synthwave",
			"electro", "brazilian bass", "hi-nrg",
		},
	},
	{
		Canonical: "Rock",
		Keywords: []string{
			"psychedelic rock", "progressive rock", "classic rock",
			"blues rock", "southern rock", "country rock", "funk rock",
			"alternative rock", "post-grunge", "garage rock", "acid rock",
			"jangle pop", "grunge", "rock and roll", "rock", "aor",
			"shoegaze", "new wave",
		},
	},
	{
		Canonical: "Pop",
		Keywords: []string{
			"soft pop", "art pop", "hyperpop", "synthpop", "britpop",
			"baroque pop", "pop punk", "pop",
		},
	},
	{
		Canonical: "R&B",
		Keywords: []string{
			"neo soul", "r&b",
		},
	},
	{
		Canonical: "Country",
		Keywords: []string{
			"country",
		},
	},
	{
		Canonical: "Punk",
		Keywords: []string{
			"proto-punk", "post-punk", "punk",
		},
	},
	{
		Canonical: "Classical",
		Keywords: []string{
			"classical", "opera", "score", "soundtrack", "cinematic",
		},
	},
	{
		Canonical: "Latin",
		Keywords: []string{
			"reggaeton", "urbano latino", "latin",
		},
	},
	{
		Canonical: "Jazz",
		Keywords: []string{
			"jazz rap", "jazz",
		},
	},
	{
		Canonical: "Blues",
		Keywords: []string{
			"blues",
		},
	},
	{
		Canonical: "Soul",
		Keywords: []string{
			"gospel", "soul",
		},
	},
	{
		Canonical: "Funk",
		Keywords: []string{
			"funk",
		},
	},
	{
		Canonical: "Reggae",
		Keywords: []string{
			"reggae", "dancehall",
		},
	},
}

// Resolves genre.
func ResolveGenre(spotifyGenres []string) string {
	if len(spotifyGenres) == 0 {
		return ""
	}

	genre := strings.ToLower(spotifyGenres[0])

	for _, mapping := range genreMapping {
		for _, keyword := range mapping.Keywords {
			if strings.Contains(genre, keyword) {
				return mapping.Canonical
			}
		}
	}

	return ""
}

// Normalizes aigenre label.
func NormalizeAIGenreLabel(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = strings.ReplaceAll(s, "hip-hop", "hip hop")
	s = strings.ReplaceAll(s, "–", "-")
	return strings.TrimSpace(s)
}

// Resolves single genre.
func ResolveSingleGenre(spotifyGenre string) string {
	g := strings.ToLower(strings.TrimSpace(spotifyGenre))
	if g == "" {
		return ""
	}
	for _, mapping := range genreMapping {
		for _, keyword := range mapping.Keywords {
			if strings.Contains(g, keyword) {
				return mapping.Canonical
			}
		}
	}
	return ""
}

// Matches genre filter.
func MatchesGenreFilter(spotifyGenres []string, resolvedAggregate, requestedCanonical string) bool {
	if requestedCanonical == "" {
		return true
	}
	for _, g := range spotifyGenres {
		if r := ResolveSingleGenre(g); r != "" && strings.EqualFold(r, requestedCanonical) {
			return true
		}
	}
	return strings.EqualFold(resolvedAggregate, requestedCanonical)
}
