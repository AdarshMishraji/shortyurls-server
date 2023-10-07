package database

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type UserLoginHistory struct {
	gorm.Model
	ID        *uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:",omitempty"`
	UserID    *uuid.UUID     `gorm:"type:uuid;not null" json:"user_id" json:",omitempty"`
	User      *User          `json:",omitempty"`
	Location  datatypes.JSON `json:"location" gorm:"type:jsonb;default:null"`
	Device    datatypes.JSON `json:"device" gorm:"type:jsonb;default:null"`
	CreatedAt *time.Time     `gorm:"default:null" json:",omitempty"`
	UpdatedAt *time.Time     `gorm:"default:null" json:",omitempty"`
	DeletedAt *time.Time     `gorm:"default:null" json:",omitempty"`
}
