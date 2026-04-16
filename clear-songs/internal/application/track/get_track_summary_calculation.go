package track

import (
	"context"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/RubenPari/clear-songs/internal/domain/shared/utils"
	domainTrack "github.com/RubenPari/clear-songs/internal/domain/track"
	spotifyAPI "github.com/zmb3/spotify"
)

// artistData holds artist information for summary calculation.
type artistData struct {
	id    string
	name  string
	count int
}

// calculateSummary calculates artist summary from tracks.
func (uc *GetTrackSummaryUseCase) calculateSummary(
	ctx context.Context,
	tracks []spotifyAPI.SavedTrack,
	min, max int,
	genre string,
) []domainTrack.ArtistSummary {
	artistMap := groupTracksByPrimaryArtist(tracks)
	artistIDs := collectArtistIDs(artistMap)
	artistDetails := uc.fetchArtistDetails(ctx, artistIDs)

	return uc.buildArtistSummary(ctx, artistMap, artistDetails, min, max, genre)
}

// groupTracksByPrimaryArtist groups tracks by their primary artist.
func groupTracksByPrimaryArtist(tracks []spotifyAPI.SavedTrack) map[string]artistData {
	artistMap := make(map[string]artistData)

	for _, savedTrack := range tracks {
		if len(savedTrack.Artists) == 0 {
			continue
		}

		artist := savedTrack.Artists[0]
		artistID := string(artist.ID)
		if artistID == "" {
			artistID = strings.ToLower(strings.TrimSpace(artist.Name))
		}

		existing := artistMap[artistID]
		existing.id = string(artist.ID)
		existing.name = artist.Name
		existing.count++
		artistMap[artistID] = existing
	}

	return artistMap
}

// collectArtistIDs collects artist IDs from the artist map.
func collectArtistIDs(artistMap map[string]artistData) []spotifyAPI.ID {
	artistIDs := make([]spotifyAPI.ID, 0, len(artistMap))
	for _, data := range artistMap {
		if data.id == "" {
			continue
		}
		artistIDs = append(artistIDs, spotifyAPI.ID(data.id))
	}

	return artistIDs
}

// fetchArtistDetails fetches artist details in batches.
func (uc *GetTrackSummaryUseCase) fetchArtistDetails(ctx context.Context, artistIDs []spotifyAPI.ID) map[string]*spotifyAPI.FullArtist {
	artistDetails := make(map[string]*spotifyAPI.FullArtist)
	if len(artistIDs) == 0 {
		return artistDetails
	}

	artists, err := uc.spotifyRepo.GetArtists(ctx, artistIDs)
	if err != nil {
		log.Printf("Error batch fetching artists: %v", err)
		return artistDetails
	}

	for _, artist := range artists {
		if artist != nil {
			artistDetails[string(artist.ID)] = artist
		}
	}

	return artistDetails
}

func (uc *GetTrackSummaryUseCase) buildArtistSummary(
	ctx context.Context,
	artistMap map[string]artistData,
	artistDetails map[string]*spotifyAPI.FullArtist,
	min, max int,
	genre string,
) []domainTrack.ArtistSummary {
	summary := make([]domainTrack.ArtistSummary, 0, len(artistMap))

	for _, data := range artistMap {
		if !passesRangeFilter(data.count, min, max) {
			continue
		}

		imageURL, genres := extractArtistMetadata(data.id, artistDetails)
		resolvedGenre := uc.resolveGenreWithFallback(ctx, data.name, genres)
		if !domainTrack.MatchesGenreFilter(genres, resolvedGenre, genre) {
			continue
		}

		if genres == nil {
			genres = []string{}
		}

		summary = append(summary, domainTrack.ArtistSummary{
			ID:       data.id,
			Name:     data.name,
			Count:    data.count,
			ImageURL: imageURL,
			Genres:   genres,
			Genre:    resolvedGenre,
		})
	}

	return summary
}

// passesRangeFilter checks if the count of tracks by an artist falls within the specified range.
func passesRangeFilter(count, min, max int) bool {
	if min > 0 && count < min {
		return false
	}
	if max > 0 && count > max {
		return false
	}

	return true
}

// extractArtistMetadata extracts the image URL and genres from the artist details.
func extractArtistMetadata(artistID string, artistDetails map[string]*spotifyAPI.FullArtist) (string, []string) {
	if artist, ok := artistDetails[artistID]; ok {
		return utils.GetMediumImage(artist.Images), artist.Genres
	}

	return "", nil
}

// resolveGenreWithFallback attempts to resolve the genre using the primary method, and falls back to AI if needed.
func (uc *GetTrackSummaryUseCase) resolveGenreWithFallback(ctx context.Context, artistName string, genres []string) string {
	resolvedGenre := domainTrack.ResolveGenre(genres)
	if resolvedGenre != "" {
		return resolvedGenre
	}
	if uc.aiRepo == nil {
		return ""
	}

	aiCtx, cancel := context.WithTimeout(ctx, geminiRequestTimeout())
	defer cancel()

	log.Printf("[genre] AI fallback: calling API artist=%q spotifyGenres=%v", artistName, genres)

	aiGenre, err := uc.aiRepo.ResolveArtistGenre(aiCtx, artistName)
	if err != nil {
		log.Printf("[genre] AI fallback: ERROR artist=%q err=%v", artistName, err)
		return ""
	}

	if aiGenre == "" {
		log.Printf("[genre] AI fallback: empty response artist=%q", artistName)
		return ""
	}

	normalized := domainTrack.NormalizeAIGenreLabel(aiGenre)
	canonical := domainTrack.ResolveGenre([]string{normalized})
	if canonical == "" {
		log.Printf("[genre] AI fallback: UNMAPPED artist=%q aiRaw=%q (no keyword matched canonical mapping)", artistName, aiGenre)
		return ""
	}
	log.Printf("[genre] AI fallback: OK artist=%q aiRaw=%q canonical=%q", artistName, aiGenre, canonical)
	return canonical
}

// geminiRequestTimeout bounds each Gemini call (default 25s; env GEMINI_REQUEST_TIMEOUT_SEC 5–120).
func geminiRequestTimeout() time.Duration {
	const defaultSec = 25
	s := strings.TrimSpace(os.Getenv("GEMINI_REQUEST_TIMEOUT_SEC"))
	if s == "" {
		return defaultSec * time.Second
	}
	sec, err := strconv.Atoi(s)
	if err != nil || sec < 5 || sec > 120 {
		return defaultSec * time.Second
	}
	return time.Duration(sec) * time.Second
}

