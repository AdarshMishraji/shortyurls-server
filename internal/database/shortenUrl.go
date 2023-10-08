package database

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ShortenURL struct {
	gorm.Model
	ID             *uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:",omitempty"`
	OriginalURL    *string        `json:"original_url,omitempty" gorm:"unique;not null"`
	Alias          *string        `json:"alias,omitempty" gorm:"unique;not null"`
	UserID         *uuid.UUID     `gorm:"type:uuid;not null" json:"user_id,omitempty"`
	IsActive       bool           `gorm:"default:true" json:"is_active"`
	IsDeleted      bool           `gorm:"default:false" json:"is_deleted"`
	User           *User          `json:",omitempty"`
	Meta           datatypes.JSON `json:"meta,omitempty" gorm:"type:jsonb;default:null"`
	ExpirationTime *time.Time     `json:"expiration_time,omitempty" gorm:"default:null"`
	Password       *string        `json:"password,omitempty" gorm:"default:null"`
}
