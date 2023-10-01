package database

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ShortenURLVisit struct {
	gorm.Model
	ID           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id" gorm:"primaryKey"`
	ShortenURLID uuid.UUID `gorm:"type:uuid" json:"shorten_url_id"`
	ShortenURL   ShortenURL
	Location     JSON `json:"location"`
	Device       JSON `json:"device"`
}
