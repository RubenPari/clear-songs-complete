package models

import (
	"gorm.io/gorm"
)

type TrackDB struct {
	gorm.Model
	Id     string `gorm:"type:varchar(100);not null"`
	Name   string `gorm:"type:varchar(100);not null"`
	Artist string `gorm:"type:varchar(100);not null"`
	Album  string `gorm:"type:varchar(100);not null"`
	URI    string `gorm:"type:varchar(200);not null"`
	URL    string `gorm:"type:varchar(200);not null"`
}
