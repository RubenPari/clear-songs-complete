package playlist

import (
	"context"

	"github.com/RubenPari/clear-songs/internal/domain/shared"
	spotifyAPI "github.com/zmb3/spotify"
)

// DeletePlaylistTracksUseCase handles the business logic for deleting tracks from a playlist
type DeletePlaylistTracksUseCase struct {
	spotifyRepo shared.SpotifyRepository
	cacheRepo   shared.CacheRepository
}

// Creates delete playlist tracks use case.
func NewDeletePlaylistTracksUseCase(
	spotifyRepo shared.SpotifyRepository,
	cacheRepo shared.CacheRepository,
) *DeletePlaylistTracksUseCase {
	return &DeletePlaylistTracksUseCase{
		spotifyRepo: spotifyRepo,
		cacheRepo:   cacheRepo,
	}
}

// Execute.
func (uc *DeletePlaylistTracksUseCase) Execute(ctx context.Context, playlistID spotifyAPI.ID) error {
	tracks, err := fetchPlaylistTracks(ctx, uc.spotifyRepo, uc.cacheRepo, playlistID)
	if err != nil {
		return err
	}

	if len(tracks) == 0 {
		return nil
	}

	trackIDs := extractPlaylistTrackIDs(tracks)
	if err := uc.spotifyRepo.DeletePlaylistTracks(ctx, playlistID, trackIDs); err != nil {
		return err
	}

	_ = uc.cacheRepo.InvalidatePlaylistTracks(ctx, playlistID)
	return nil
}
