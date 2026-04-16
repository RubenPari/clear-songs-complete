package track

// TrackResponse represents a track in API responses
type TrackResponse struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Artists    []string `json:"artists"`
	Album      string   `json:"album"`
	Duration   int      `json:"duration"`
	ImageURL   string   `json:"image_url,omitempty"`
	SpotifyURL string   `json:"spotify_url,omitempty"`
}

// RangeRequest binds optional min/max query params. Absent keys stay nil so "min only"
// does not force max=0 and break validation (see ValidateRangeQuery).
type RangeRequest struct {
	Min   *int   `form:"min" binding:"omitempty,min=0"`
	Max   *int   `form:"max" binding:"omitempty,min=0"`
	Genre string `form:"genre"`
}

// Validate range query.
func ValidateRangeQuery(req *RangeRequest) (min, max int, errMsg string) {
	if req == nil {
		return 0, 0, ""
	}
	if req.Min != nil {
		if *req.Min < 0 {
			return 0, 0, "min must be >= 0"
		}
		min = *req.Min
	}
	if req.Max != nil {
		if *req.Max < 0 {
			return 0, 0, "max must be >= 0"
		}
		max = *req.Max
	}
	if max > 0 && min > max {
		return 0, 0, "min must be less than or equal to max"
	}
	return min, max, ""
}
