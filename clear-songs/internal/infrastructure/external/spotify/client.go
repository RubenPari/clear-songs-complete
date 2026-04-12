package spotify

import (
	"context"
	"errors"
	"log"
	"sync"

	"github.com/RubenPari/clear-songs/internal/domain/shared"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

// SpotifyRepositoryImpl implements the SpotifyRepository interface.
type SpotifyRepositoryImpl struct {
	authenticator spotify.Authenticator
	mu            sync.RWMutex
	client        *spotify.Client
}

// NewSpotifyRepository creates a new Spotify repository implementation.
func NewSpotifyRepository(clientID, clientSecret, redirectURI string, scopes []string) *SpotifyRepositoryImpl {
	auth := spotify.NewAuthenticator(redirectURI, scopes...)
	auth.SetAuthInfo(clientID, clientSecret)

	return &SpotifyRepositoryImpl{
		authenticator: auth,
	}
}

// SetAccessToken sets the OAuth token and creates a new client.
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

// GetClient returns the Spotify client (for backward compatibility).
func (r *SpotifyRepositoryImpl) GetClient() *spotify.Client {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.client
}

func (r *SpotifyRepositoryImpl) currentClient() (*spotify.Client, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.client == nil {
		return nil, errors.New("spotify client not initialized")
	}

	return r.client, nil
}

// GetCurrentUser retrieves the current authenticated user.
func (r *SpotifyRepositoryImpl) GetCurrentUser(ctx context.Context) (*spotify.PrivateUser, error) {
	client, err := r.currentClient()
	if err != nil {
		return nil, err
	}
	return client.CurrentUser()
}

// GetUserTracks retrieves tracks saved by the user with pagination.
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

// GetAllUserTracks retrieves all user tracks with automatic pagination.
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

	log.Printf("Retrieved %d total tracks", len(allTracks))
	return allTracks, nil
}

// GetTracksByArtist filters tracks by artist ID.
func (r *SpotifyRepositoryImpl) GetTracksByArtist(ctx context.Context, artistID spotify.ID, tracks []spotify.SavedTrack) ([]spotify.SavedTrack, error) {
	var filteredTracks []spotify.SavedTrack

	for _, track := range tracks {
		if len(track.Artists) > 0 && track.Artists[0].ID == artistID {
			filteredTracks = append(filteredTracks, track)
		}
	}

	return filteredTracks, nil
}

// GetTrackIDsByArtist returns only track IDs for an artist.
func (r *SpotifyRepositoryImpl) GetTrackIDsByArtist(ctx context.Context, artistID spotify.ID, tracks []spotify.SavedTrack) ([]spotify.ID, error) {
	var trackIDs []spotify.ID

	for _, track := range tracks {
		if len(track.Artists) > 0 && track.Artists[0].ID == artistID {
			trackIDs = append(trackIDs, track.ID)
		}
	}

	return trackIDs, nil
}

// DeleteTracksFromLibrary removes tracks from user's library.
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

		log.Printf("Deleted tracks from offset: %d", offset)
		offset += limit
	}

	return nil
}

// GetPlaylist retrieves a playlist by ID.
func (r *SpotifyRepositoryImpl) GetPlaylist(ctx context.Context, playlistID spotify.ID) (*spotify.FullPlaylist, error) {
	client, err := r.currentClient()
	if err != nil {
		return nil, err
	}
	return client.GetPlaylist(playlistID)
}

// GetPlaylistTracks retrieves tracks from a playlist with pagination.
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

// GetAllPlaylistTracks retrieves all tracks from a playlist with automatic pagination.
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

// DeletePlaylistTracks removes tracks from a playlist.
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

// GetUserPlaylists retrieves playlists owned or followed by the user with pagination.
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

// GetAllUserPlaylists retrieves all user playlists with automatic pagination.
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

// GetArtist retrieves artist information.
func (r *SpotifyRepositoryImpl) GetArtist(ctx context.Context, artistID spotify.ID) (*spotify.FullArtist, error) {
	client, err := r.currentClient()
	if err != nil {
		return nil, err
	}
	return client.GetArtist(artistID)
}

// GetArtists retrieves multiple artists in batch (up to 50 per Spotify API call).
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
			log.Printf("Error fetching artists batch [%d:%d]: %v", i, end, err)
			return nil, err
		}

		allArtists = append(allArtists, batch...)
	}

	return allArtists, nil
}

// GetTrack retrieves track information.
func (r *SpotifyRepositoryImpl) GetTrack(ctx context.Context, trackID spotify.ID) (*spotify.FullTrack, error) {
	client, err := r.currentClient()
	if err != nil {
		return nil, err
	}
	return client.GetTrack(trackID)
}

// Ensure SpotifyRepositoryImpl implements SpotifyRepository interface.
var _ shared.SpotifyRepository = (*SpotifyRepositoryImpl)(nil)
