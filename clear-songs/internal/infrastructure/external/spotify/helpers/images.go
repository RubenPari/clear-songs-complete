// Package helpers provides utility functions for working with Spotify API data,
// such as selecting appropriately sized images from Spotify's image sets.
package helpers

import spotifyAPI "github.com/zmb3/spotify"

// GetSmallestImage returns the smallest image URL that fits within maxWidth.
// If no image fits, it falls back to the smallest available image.
func GetSmallestImage(images []spotifyAPI.Image, maxWidth int) string {
	if len(images) == 0 {
		return ""
	}

	for i := len(images) - 1; i >= 0; i-- {
		if images[i].Width <= maxWidth || i == 0 {
			return images[i].URL
		}
	}

	return ""
}

// GetMediumImage returns a medium-sized image URL from the provided images.
func GetMediumImage(images []spotifyAPI.Image) string {
	return GetSmallestImage(images, 300)
}
