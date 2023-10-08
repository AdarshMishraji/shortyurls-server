package database

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type UserLoginHistory struct {
	gorm.Model
	ID       *uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:",omitempty"`
	UserID   *uuid.UUID     `gorm:"type:uuid;not null" json:"user_id,omitempty"`
	User     *User          `json:",omitempty"`
	Location datatypes.JSON `json:"location" gorm:"type:jsonb;default:null"`
	Device   datatypes.JSON `json:"device" gorm:"type:jsonb;default:null"`
}
