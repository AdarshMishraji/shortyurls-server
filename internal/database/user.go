package database

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID        *uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:",omitempty"`
	Name      *string    `json:"name" gorm:"not null" json:",omitempty"`
	Email     *string    `gorm:"unique;not null" json:"email" json:",omitempty"`
	Picture   *string    `json:"picture" json:",omitempty"`
	Provider  *string    `json:"provider" json:",omitempty"`
	IsDeleted bool       `gorm:"default:false" json:"is_deleted"`
	CreatedAt *time.Time `gorm:"default:null" json:",omitempty"`
	UpdatedAt *time.Time `gorm:"default:null" json:",omitempty"`
	DeletedAt *time.Time `gorm:"default:null" json:",omitempty"`
}
