package utils

import "github.com/zmb3/spotify"

// GetSmallestImage returns the smallest image URL that fits within maxWidth,
// or the smallest available image if none fit
func GetSmallestImage(images []spotify.Image, maxWidth int) string {
	if len(images) == 0 {
		return ""
	}
	
	// Images are typically sorted from largest to smallest
	// Iterate from the end to find the smallest that fits
	for i := len(images) - 1; i >= 0; i-- {
		if images[i].Width <= maxWidth || i == 0 {
			return images[i].URL
		}
	}
	
	return ""
}

// GetMediumImage returns a medium-sized image (around 300-640px width)
func GetMediumImage(images []spotify.Image) string {
	return GetSmallestImage(images, 300)
}
