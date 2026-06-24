package shared

import "fmt"

// Cache key constants used across cache implementations and application use cases.
const (
	CacheKeyTokenPrefix          = "spotify_token"
	CacheKeyUserTracks           = "userTracks"
	CacheKeyUserPlaylists        = "userPlaylists"
	CacheKeyPlaylistTracksFmt    = "tracksPlaylist%s"
	CacheKeyTrackSummaryFmt      = "track_summary_%d_%d"
	CacheKeyTrackSummaryGenreFmt = "track_summary_%d_%d_%s"
	CacheKeyArtistAIGenreFmt     = "artist_ai_genre:%s"
)

// CacheKeyPlaylistTracks returns the Redis key for a playlist's tracks.
func CacheKeyPlaylistTracks(playlistID string) string {
	return fmt.Sprintf(CacheKeyPlaylistTracksFmt, playlistID)
}
