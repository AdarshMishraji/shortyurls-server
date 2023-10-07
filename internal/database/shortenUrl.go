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
	OriginalURL    *string        `json:"original_url" gorm:"unique;not null" json:",omitempty"`
	Alias          *string        `json:"alias" gorm:"unique;not null" json:",omitempty"`
	UserID         *uuid.UUID     `gorm:"type:uuid;not null" json:"user_id" json:",omitempty"`
	IsActive       bool           `gorm:"default:true" json:"is_active"`
	IsDeleted      bool           `gorm:"default:false" json:"is_deleted"`
	User           *User          `json:",omitempty"`
	Meta           datatypes.JSON `json:"meta" gorm:"type:jsonb" gorm:"default:null" json:",omitempty"`
	ExpirationTime *time.Time     `json:"expiration_time" gorm:"default:null" json:",omitempty"`
	Password       *string        `json:"password" json:",omitempty" gorm:"default:null"`
	CreatedAt      *time.Time     `gorm:"default:null" json:",omitempty"`
	UpdatedAt      *time.Time     `gorm:"default:null" json:",omitempty"`
	DeletedAt      *time.Time     `gorm:"default:null" json:",omitempty"`
}
