package database

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ShortenURL struct {
	gorm.Model
	ID             uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id" gorm:"primaryKey"`
	OriginalURL    string    `json:"original_url" gorm:"uniqueIndex"`
	Alias          string    `json:"alias" gorm:"uniqueIndex"`
	UserID         uuid.UUID `gorm:"type:uuid" json:"user_id"`
	IsActive       bool      `gorm:"default:true" json:"is_active"`
	IsDeleted      bool      `gorm:"default:false" json:"is_deleted"`
	User           User
	Meta           JSON      `json:"meta"`
	ExpirationTime time.Time `json:"expiration_time" gorm:"default:null"`
	Password       string    `json:"password"`
}
