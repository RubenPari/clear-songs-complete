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
			"hip hop", "rap", "trap", "boom bap", "gangster",
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
			"classical", "opera", "score", "soundtrack",
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
			"soul",
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

// ResolveGenre maps a Spotify genres array to a single canonical genre.
// It uses only the first element of the array (priority genre).
// Returns "" if no match is found or the array is empty.
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
