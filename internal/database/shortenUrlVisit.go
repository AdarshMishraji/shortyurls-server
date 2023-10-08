package database

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ShortenURLVisit struct {
	gorm.Model
	ID           *uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:",omitempty"`
	ShortenURLID *uuid.UUID     `gorm:"type:uuid;not null" json:"shorten_url_id,omitempty"`
	ShortenURL   *ShortenURL    `json:",omitempty"`
	Location     datatypes.JSON `json:"location" gorm:"type:jsonb;default:null"`
	Device       datatypes.JSON `json:"device" gorm:"type:jsonb;default:null"`
}
