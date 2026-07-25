// Package models defines GORM-backed persistence models for the Postgres
// database layer. These structs map directly to database tables and are
// used by the Postgres repository for backup and restore operations.
package models

import (
	"gorm.io/gorm"
)

// TrackDB is the GORM model representing a single track stored in the
// Postgres backup table. It mirrors the essential fields of a Spotify
// track so that library data can be persisted independently of the
// Spotify API.
type TrackDB struct {
	gorm.Model
	Id     string `gorm:"type:varchar(100);not null"`
	Name   string `gorm:"type:varchar(100);not null"`
	Artist string `gorm:"type:varchar(100);not null"`
	Album  string `gorm:"type:varchar(100);not null"`
	URI    string `gorm:"type:varchar(200);not null"`
	URL    string `gorm:"type:varchar(200);not null"`
}
