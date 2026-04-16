package track

import (
	"context"
	"strings"

	"github.com/RubenPari/clear-songs/internal/domain/shared"
	"github.com/RubenPari/clear-songs/internal/domain/shared/utils"
	domainTrack "github.com/RubenPari/clear-songs/internal/domain/track"
	spotifyAPI "github.com/zmb3/spotify"
	"go.uber.org/zap"
)

// artistData holds artist information for summary calculation.
type artistData struct {
	id    string
	name  string
	count int
}

// Calculate summary.
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

// Group tracks by primary artist.
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

// Collect artist ids.
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

// Fetch artist details.
func (uc *GetTrackSummaryUseCase) fetchArtistDetails(ctx context.Context, artistIDs []spotifyAPI.ID) map[string]*spotifyAPI.FullArtist {
	artistDetails := make(map[string]*spotifyAPI.FullArtist)
	if len(artistIDs) == 0 {
		return artistDetails
	}

	artists, err := uc.spotifyRepo.GetArtists(ctx, artistIDs)
	if err != nil {
		zap.L().Warn("error batch fetching artists", zap.Error(err))
		return artistDetails
	}

	for _, artist := range artists {
		if artist != nil {
			artistDetails[string(artist.ID)] = artist
		}
	}

	return artistDetails
}

// Builds artist summary.
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
		zap.L().Info("resolving AI genres", zap.Int("artist_count", len(needsAI)))
		rawMap, err := uc.aiRepo.ResolveArtistGenres(ctx, needsAI)
		if err != nil {
			zap.L().Warn("AI genre batch failed", zap.Error(err))
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

// Applies aigenre result.
func (uc *GetTrackSummaryUseCase) applyAIGenreResult(ctx context.Context, artistName, mapKey, aiRaw string) string {
	if aiRaw == "" {
		zap.L().Debug("AI fallback empty genre", zap.String("artist", artistName))
		return ""
	}
	normalized := domainTrack.NormalizeAIGenreLabel(aiRaw)
	canonical := domainTrack.ResolveGenre([]string{normalized})
	if canonical == "" {
		zap.L().Debug("AI fallback unmapped genre", zap.String("artist", artistName), zap.String("ai_raw", aiRaw))
		return ""
	}
	zap.L().Debug("AI fallback resolved genre",
		zap.String("artist", artistName),
		zap.String("ai_raw", aiRaw),
		zap.String("canonical", canonical),
	)
	uc.setCachedArtistCanonicalGenre(ctx, mapKey, canonical)
	return canonical
}

// Passes range filter.
func passesRangeFilter(count, min, max int) bool {
	if min > 0 && count < min {
		return false
	}
	if max > 0 && count > max {
		return false
	}

	return true
}

// Extract artist metadata.
func extractArtistMetadata(artistID string, artistDetails map[string]*spotifyAPI.FullArtist) (string, []string) {
	if artist, ok := artistDetails[artistID]; ok {
		return utils.GetMediumImage(artist.Images), artist.Genres
	}

	return "", nil
}
