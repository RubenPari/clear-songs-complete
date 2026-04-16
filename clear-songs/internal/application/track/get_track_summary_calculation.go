package track

import (
	"context"
	"log"
	"strings"

	"github.com/RubenPari/clear-songs/internal/domain/shared"
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
	resolvedByKey := make(map[string]string)
	var needsAI []shared.AIArtistLookup

	for mapKey, data := range artistMap {
		if !passesRangeFilter(data.count, min, max) {
			continue
		}

		_, genres := extractArtistMetadata(data.id, artistDetails)
		if g := domainTrack.ResolveGenre(genres); g != "" {
			resolvedByKey[mapKey] = g
			continue
		}

		if cached, ok := uc.getCachedArtistCanonicalGenre(ctx, mapKey); ok {
			resolvedByKey[mapKey] = cached
			continue
		}

		needsAI = append(needsAI, shared.AIArtistLookup{Key: mapKey, Name: data.name})
	}

	if len(needsAI) > 0 && uc.aiRepo != nil {
		log.Printf("[genre] AI batch: resolving %d artist(s) without Spotify-mappable genres", len(needsAI))
		rawMap, err := uc.aiRepo.ResolveArtistGenres(ctx, needsAI)
		if err != nil {
			log.Printf("[genre] AI batch: ERROR %v — artists left unresolved", err)
			for _, l := range needsAI {
				resolvedByKey[l.Key] = ""
			}
		} else {
			for _, l := range needsAI {
				aiRaw := ""
				if rawMap != nil {
					aiRaw = rawMap[l.Key]
				}
				resolvedByKey[l.Key] = uc.applyAIGenreResult(ctx, l.Name, l.Key, aiRaw)
			}
		}
	}

	summary := make([]domainTrack.ArtistSummary, 0, len(artistMap))

	for mapKey, data := range artistMap {
		if !passesRangeFilter(data.count, min, max) {
			continue
		}

		imageURL, genres := extractArtistMetadata(data.id, artistDetails)
		resolvedGenre := resolvedByKey[mapKey]

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

func (uc *GetTrackSummaryUseCase) applyAIGenreResult(ctx context.Context, artistName, mapKey, aiRaw string) string {
	if aiRaw == "" {
		log.Printf("[genre] AI fallback: empty RAW artist=%q", artistName)
		return ""
	}
	normalized := domainTrack.NormalizeAIGenreLabel(aiRaw)
	canonical := domainTrack.ResolveGenre([]string{normalized})
	if canonical == "" {
		log.Printf("[genre] AI fallback: UNMAPPED artist=%q aiRaw=%q (no keyword matched canonical mapping)", artistName, aiRaw)
		return ""
	}
	log.Printf("[genre] AI fallback: OK artist=%q aiRaw=%q canonical=%q", artistName, aiRaw, canonical)
	uc.setCachedArtistCanonicalGenre(ctx, mapKey, canonical)
	return canonical
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
