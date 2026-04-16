package postgres

import (
	"errors"

	"github.com/RubenPari/clear-songs/internal/domain/shared"
	"github.com/RubenPari/clear-songs/internal/infrastructure/persistence/postgres/models"
	spotifyAPI "github.com/zmb3/spotify"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// PostgresRepository implements DatabaseRepository interface
type PostgresRepository struct {
	db *gorm.DB
}

// First artist name.
func firstArtistName(artists []spotifyAPI.SimpleArtist) string {
	if len(artists) == 0 {
		return "Unknown Artist"
	}
	if artists[0].Name == "" {
		return "Unknown Artist"
	}
	return artists[0].Name
}

// Spotify url.
func spotifyURL(urls map[string]string) string {
	if urls == nil {
		return ""
	}
	return urls["spotify"]
}

// Creates postgres repository.
func NewPostgresRepository(db *gorm.DB) shared.DatabaseRepository {
	if db == nil {
		return &NoOpDatabaseRepository{}
	}
	return &PostgresRepository{db: db}
}

// Save tracks backup.
func (r *PostgresRepository) SaveTracksBackup(tracks []spotifyAPI.PlaylistTrack) error {
	zap.L().Info("saving tracks backup started", zap.Int("count", len(tracks)))

	for _, trackPlaylist := range tracks {
		trackID := trackPlaylist.Track.ID.String()
		if trackID == "" {
			continue
		}
		track := models.TrackDB{
			Id:     trackID,
			Name:   trackPlaylist.Track.Name,
			Artist: firstArtistName(trackPlaylist.Track.Artists),
			Album:  trackPlaylist.Track.Album.Name,
			URI:    string(trackPlaylist.Track.URI),
			URL:    spotifyURL(trackPlaylist.Track.ExternalURLs),
		}

		if err := r.saveToDB(track); err != nil {
			return err
		}
	}

	return nil
}

// Save full tracks backup.
func (r *PostgresRepository) SaveFullTracksBackup(tracks []spotifyAPI.FullTrack) error {
	zap.L().Info("saving full tracks backup started", zap.Int("count", len(tracks)))

	for _, t := range tracks {
		trackID := t.ID.String()
		if trackID == "" {
			continue
		}
		track := models.TrackDB{
			Id:     trackID,
			Name:   t.Name,
			Artist: firstArtistName(t.Artists),
			Album:  t.Album.Name,
			URI:    string(t.URI),
			URL:    spotifyURL(t.ExternalURLs),
		}

		if err := r.saveToDB(track); err != nil {
			return err
		}
	}

	return nil
}

// Save to db.
func (r *PostgresRepository) saveToDB(track models.TrackDB) error {
	var existingTrack models.TrackDB
	result := r.db.First(&existingTrack, "id = ?", track.Id)

	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			zap.L().Error("error querying existing track", zap.Error(result.Error))
			return result.Error
		}

		// Track doesn't exist, insert it
		if err := r.db.Create(&track).Error; err != nil {
			zap.L().Error("error inserting track", zap.String("track_id", track.Id), zap.Error(err))
			return err
		}
	}
	return nil
}

// NoOpDatabaseRepository is a no-op implementation when database is not available
type NoOpDatabaseRepository struct{}

// Save tracks backup.
func (n *NoOpDatabaseRepository) SaveTracksBackup(tracks []spotifyAPI.PlaylistTrack) error {
	zap.L().Warn("database not available, skipping track backup")
	return nil // No-op
}

// Save full tracks backup.
func (n *NoOpDatabaseRepository) SaveFullTracksBackup(tracks []spotifyAPI.FullTrack) error {
	zap.L().Warn("database not available, skipping track backup")
	return nil // No-op
}

// Ensure implementations
var _ shared.DatabaseRepository = (*PostgresRepository)(nil)
var _ shared.DatabaseRepository = (*NoOpDatabaseRepository)(nil)
