package track

import "testing"

func TestResolveGenre(t *testing.T) {
	tests := []struct {
		name     string
		genres   []string
		expected string
	}{
		{"gangster rap -> Hip Hop", []string{"gangster rap"}, "Hip Hop"},
		{"east coast hip hop -> Hip Hop", []string{"east coast hip hop"}, "Hip Hop"},
		{"rap -> Hip Hop", []string{"rap"}, "Hip Hop"},
		{"boom bap -> Hip Hop", []string{"boom bap"}, "Hip Hop"},
		{"horrorcore -> Hip Hop", []string{"horrorcore"}, "Hip Hop"},
		{"progressive house -> Electronic", []string{"progressive house"}, "Electronic"},
		{"edm -> Electronic", []string{"edm"}, "Electronic"},
		{"trance -> Electronic", []string{"trance"}, "Electronic"},
		{"moombahton -> Electronic", []string{"moombahton"}, "Electronic"},
		{"lo-fi -> Electronic", []string{"lo-fi"}, "Electronic"},
		{"classic rock -> Rock", []string{"classic rock"}, "Rock"},
		{"grunge -> Rock", []string{"grunge"}, "Rock"},
		{"aor -> Rock", []string{"aor"}, "Rock"},
		{"southern rock -> Rock", []string{"southern rock"}, "Rock"},
		{"nu metal -> Metal", []string{"nu metal"}, "Metal"},
		{"thrash metal -> Metal", []string{"thrash metal"}, "Metal"},
		{"glam metal -> Metal", []string{"glam metal"}, "Metal"},
		{"alternative metal -> Metal", []string{"alternative metal"}, "Metal"},
		{"pop -> Pop", []string{"pop"}, "Pop"},
		{"soft pop -> Pop", []string{"soft pop"}, "Pop"},
		{"hyperpop -> Pop", []string{"hyperpop"}, "Pop"},
		{"synthpop -> Pop", []string{"synthpop"}, "Pop"},
		{"r&b -> R&B", []string{"r&b"}, "R&B"},
		{"neo soul -> R&B", []string{"neo soul"}, "R&B"},
		{"country -> Country", []string{"country"}, "Country"},
		{"punk -> Punk", []string{"punk"}, "Punk"},
		{"post-punk -> Punk", []string{"post-punk"}, "Punk"},
		{"classical -> Classical", []string{"classical"}, "Classical"},
		{"opera -> Classical", []string{"opera"}, "Classical"},
		{"reggaeton -> Latin", []string{"reggaeton"}, "Latin"},
		{"folk rock -> Folk", []string{"folk rock"}, "Folk"},
		{"folk -> Folk", []string{"folk"}, "Folk"},
		{"unknown genre -> empty", []string{"totally unknown xyz"}, ""},
		{"empty array -> empty", []string{}, ""},
		{"uses first element only", []string{"edm", "gangster rap"}, "Electronic"},
		{"rap metal -> Metal", []string{"rap metal"}, "Metal"},
		{"country hip hop -> Hip Hop", []string{"country hip hop"}, "Hip Hop"},
		{"brazilian bass -> Electronic", []string{"brazilian bass"}, "Electronic"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ResolveGenre(tt.genres)
			if result != tt.expected {
				t.Errorf("ResolveGenre(%v) = %q, want %q", tt.genres, result, tt.expected)
			}
		})
	}
}
