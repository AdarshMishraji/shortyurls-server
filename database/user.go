package database

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	Email     string    `gorm:"uniqueIndex" json:"email"`
	Picture   string    `json:"picture"`
	Provider  string    `json:"provider"`
	IsDeleted bool      `gorm:"default:false" json:"is_deleted"`
}
