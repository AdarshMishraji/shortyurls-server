package database

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID        *uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:",omitempty"`
	Name      *string    `json:"name,omitempty" gorm:"not null"`
	Email     *string    `gorm:"unique;not null" json:"email,omitempty"`
	Picture   *string    `json:"picture,omitempty"`
	Provider  *string    `json:"provider,omitempty"`
	IsDeleted bool       `gorm:"default:false" json:"is_deleted"`
}
