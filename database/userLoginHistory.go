package database

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserLoginHistory struct {
	gorm.Model
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id" gorm:"primaryKey"`
	UserID   uuid.UUID `gorm:"type:uuid" json:"user_id"`
	User     User
	Location JSON `json:"location"`
	Device   JSON `json:"device"`
}
