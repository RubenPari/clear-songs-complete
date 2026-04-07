package models

type ArtistSummary struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Count    int    `json:"count"`
	ImageURL string `json:"image_url,omitempty"`
}
