// Package spotify implements the SpotifyRepository interface using the
// zmb3/spotify client library. It handles OAuth token management, paginated
// fetching of tracks/playlists/artists, and batched mutation operations
// (delete tracks) that respect Spotify API rate limits.
package spotify

import (
	"context"
	"errors"
	"sync"

	"github.com/RubenPari/clear-songs/internal/domain/shared"
	"github.com/zmb3/spotify"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

// SpotifyRepositoryImpl implements the SpotifyRepository interface.
// It guards the underlying spotify.Client with a read-write mutex so
// that concurrent handlers can safely read the client while token
// refreshes hold an exclusive write lock.
type SpotifyRepositoryImpl struct {
	authenticator spotify.Authenticator
	mu            sync.RWMutex
	client        *spotify.Client
}

// NewSpotifyRepository creates an authenticator configured with the
// given OAuth credentials and returns an uninitialised repository.
// Call SetAccessToken before invoking any data-fetching methods.
func NewSpotifyRepository(clientID, clientSecret, redirectURI string, scopes []string) *SpotifyRepositoryImpl {
	auth := spotify.NewAuthenticator(redirectURI, scopes...)
	auth.SetAuthInfo(clientID, clientSecret)

	return &SpotifyRepositoryImpl{
		authenticator: auth,
	}
}

// SetAccessToken builds an authenticated Spotify client from the supplied
// OAuth2 token and stores it under a write lock. Passing nil clears the
// client (e.g. on logout). Returns an error when the token type assertion
// fails.
func (r *SpotifyRepositoryImpl) SetAccessToken(token interface{}) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if token == nil {
		r.client = nil
		return nil
	}

	oauthToken, ok := token.(*oauth2.Token)
	if !ok {
		return errors.New("invalid token type")
	}

	client := r.authenticator.NewClient(oauthToken)
	r.client = &client
	return nil
}

// GetClient returns the current authenticated Spotify client, or nil if
// no token has been set. The returned pointer is a snapshot; callers must
// not retain it across requests.
func (r *SpotifyRepositoryImpl) GetClient() *spotify.Client {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.client
}

// currentClient returns the authenticated client or an error when no
// token has been set. All public methods delegate here for thread-safe
// access.
func (r *SpotifyRepositoryImpl) currentClient() (*spotify.Client, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.client == nil {
		return nil, errors.New("spotify client not initialized")
	}

	return r.client, nil
}

// GetCurrentUser retrieves the Spotify profile of the authenticated user.
func (r *SpotifyRepositoryImpl) GetCurrentUser(ctx context.Context) (*spotify.PrivateUser, error) {
	client, err := r.currentClient()
	if err != nil {
		return nil, err
	}
	return client.CurrentUser()
}

// GetUserTracks fetches a single page of the authenticated user's saved
// tracks with the given limit and offset.
func (r *SpotifyRepositoryImpl) GetUserTracks(ctx context.Context, limit, offset int) ([]spotify.SavedTrack, error) {
	client, err := r.currentClient()
	if err != nil {
		return nil, err
	}

	page, err := client.CurrentUsersTracksOpt(&spotify.Options{
		Limit:  &limit,
		Offset: &offset,
	})
	if err != nil {
		return nil, err
	}

	return page.Tracks, nil
}

// GetAllUserTracks paginates through the entire saved-tracks library in
// pages of 50 and returns every track. It stops when a page returns zero
// results.
func (r *SpotifyRepositoryImpl) GetAllUserTracks(ctx context.Context) ([]spotify.SavedTrack, error) {
	var allTracks []spotify.SavedTrack
	limit := 50
	offset := 0

	for {
		tracks, err := r.GetUserTracks(ctx, limit, offset)
		if err != nil {
			return nil, err
		}

		if len(tracks) == 0 {
			break
		}

		allTracks = append(allTracks, tracks...)
		offset += limit
	}

	zap.L().Info("retrieved user tracks", zap.Int("count", len(allTracks)))
	return allTracks, nil
}

// GetTracksByArtist filters the provided tracks slice, keeping only those
// whose primary artist matches the given artistID.
func (r *SpotifyRepositoryImpl) GetTracksByArtist(ctx context.Context, artistID spotify.ID, tracks []spotify.SavedTrack) ([]spotify.SavedTrack, error) {
	var filteredTracks []spotify.SavedTrack

	for _, track := range tracks {
		if len(track.Artists) > 0 && track.Artists[0].ID == artistID {
			filteredTracks = append(filteredTracks, track)
		}
	}

	return filteredTracks, nil
}

// GetTrackIDsByArtist returns only the track IDs (not full track objects)
// for tracks whose primary artist matches the given artistID.
func (r *SpotifyRepositoryImpl) GetTrackIDsByArtist(ctx context.Context, artistID spotify.ID, tracks []spotify.SavedTrack) ([]spotify.ID, error) {
	var trackIDs []spotify.ID

	for _, track := range tracks {
		if len(track.Artists) > 0 && track.Artists[0].ID == artistID {
			trackIDs = append(trackIDs, track.ID)
		}
	}

	return trackIDs, nil
}

// DeleteTracksFromLibrary removes the given track IDs from the user's
// saved library in batches of 50 to stay within Spotify API limits.
func (r *SpotifyRepositoryImpl) DeleteTracksFromLibrary(ctx context.Context, trackIDs []spotify.ID) error {
	client, err := r.currentClient()
	if err != nil {
		return err
	}

	limit := 50
	offset := 0

	for offset < len(trackIDs) {
		end := offset + limit
		if end > len(trackIDs) {
			end = len(trackIDs)
		}

		batch := trackIDs[offset:end]
		if err := client.RemoveTracksFromLibrary(batch...); err != nil {
			return err
		}

		zap.L().Info("deleted track batch", zap.Int("offset", offset), zap.Int("batch_size", len(batch)))
		offset += limit
	}

	return nil
}

// GetPlaylist fetches full metadata for a single playlist by ID.
func (r *SpotifyRepositoryImpl) GetPlaylist(ctx context.Context, playlistID spotify.ID) (*spotify.FullPlaylist, error) {
	client, err := r.currentClient()
	if err != nil {
		return nil, err
	}
	return client.GetPlaylist(playlistID)
}

// GetPlaylistTracks fetches a single page of tracks from the specified
// playlist with the given limit and offset.
func (r *SpotifyRepositoryImpl) GetPlaylistTracks(ctx context.Context, playlistID spotify.ID, limit, offset int) ([]spotify.PlaylistTrack, error) {
	client, err := r.currentClient()
	if err != nil {
		return nil, err
	}

	page, err := client.GetPlaylistTracksOpt(playlistID, &spotify.Options{
		Offset: &offset,
		Limit:  &limit,
	}, "")
	if err != nil {
		return nil, err
	}

	return page.Tracks, nil
}

// GetAllPlaylistTracks paginates through every track in a playlist in
// pages of 100 and returns the complete list.
func (r *SpotifyRepositoryImpl) GetAllPlaylistTracks(ctx context.Context, playlistID spotify.ID) ([]spotify.PlaylistTrack, error) {
	limit := 100
	offset := 0
	var allTracks []spotify.PlaylistTrack

	for {
		tracks, err := r.GetPlaylistTracks(ctx, playlistID, limit, offset)
		if err != nil {
			return nil, err
		}

		if len(tracks) < limit {
			allTracks = append(allTracks, tracks...)
			break
		}

		allTracks = append(allTracks, tracks...)
		offset += limit
	}

	return allTracks, nil
}

// DeletePlaylistTracks removes the specified tracks from a playlist in
// batches of 100 to comply with Spotify API constraints.
func (r *SpotifyRepositoryImpl) DeletePlaylistTracks(ctx context.Context, playlistID spotify.ID, trackIDs []spotify.ID) error {
	client, err := r.currentClient()
	if err != nil {
		return err
	}

	limit := 100
	offset := 0

	for offset < len(trackIDs) {
		end := offset + limit
		if end > len(trackIDs) {
			end = len(trackIDs)
		}

		batch := trackIDs[offset:end]
		if _, err := client.RemoveTracksFromPlaylist(playlistID, batch...); err != nil {
			return err
		}

		offset += limit
	}

	return nil
}

// GetUserPlaylists fetches a single page of the authenticated user's
// playlists with the given limit and offset.
func (r *SpotifyRepositoryImpl) GetUserPlaylists(ctx context.Context, limit, offset int) ([]spotify.SimplePlaylist, error) {
	client, err := r.currentClient()
	if err != nil {
		return nil, err
	}

	page, err := client.CurrentUsersPlaylistsOpt(&spotify.Options{
		Limit:  &limit,
		Offset: &offset,
	})
	if err != nil {
		return nil, err
	}

	return page.Playlists, nil
}

// GetAllUserPlaylists paginates through all of the authenticated user's
// playlists in pages of 50 and returns the complete list.
func (r *SpotifyRepositoryImpl) GetAllUserPlaylists(ctx context.Context) ([]spotify.SimplePlaylist, error) {
	var allPlaylists []spotify.SimplePlaylist
	limit := 50
	offset := 0

	for {
		playlists, err := r.GetUserPlaylists(ctx, limit, offset)
		if err != nil {
			return nil, err
		}

		if len(playlists) == 0 {
			break
		}

		allPlaylists = append(allPlaylists, playlists...)
		offset += limit
	}

	return allPlaylists, nil
}

// GetArtist fetches full metadata for a single artist by ID.
func (r *SpotifyRepositoryImpl) GetArtist(ctx context.Context, artistID spotify.ID) (*spotify.FullArtist, error) {
	client, err := r.currentClient()
	if err != nil {
		return nil, err
	}
	return client.GetArtist(artistID)
}

// GetArtists fetches full metadata for multiple artists in batches of 50,
// respecting the Spotify API batch limit for the GetArtists endpoint.
func (r *SpotifyRepositoryImpl) GetArtists(ctx context.Context, artistIDs []spotify.ID) ([]*spotify.FullArtist, error) {
	client, err := r.currentClient()
	if err != nil {
		return nil, err
	}

	var allArtists []*spotify.FullArtist
	batchSize := 50

	for i := 0; i < len(artistIDs); i += batchSize {
		end := i + batchSize
		if end > len(artistIDs) {
			end = len(artistIDs)
		}

		batch, err := client.GetArtists(artistIDs[i:end]...)
		if err != nil {
			zap.L().Error("error fetching artists batch",
				zap.Int("start", i),
				zap.Int("end", end),
				zap.Error(err),
			)
			return nil, err
		}

		allArtists = append(allArtists, batch...)
	}

	return allArtists, nil
}

// GetTrack fetches full metadata for a single track by ID.
func (r *SpotifyRepositoryImpl) GetTrack(ctx context.Context, trackID spotify.ID) (*spotify.FullTrack, error) {
	client, err := r.currentClient()
	if err != nil {
		return nil, err
	}
	return client.GetTrack(trackID)
}

// Ensure SpotifyRepositoryImpl implements SpotifyRepository interface.
var _ shared.SpotifyRepository = (*SpotifyRepositoryImpl)(nil)
