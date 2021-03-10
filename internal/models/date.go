package models

import (
	"time"

	"gorm.io/gorm"
)

type Date struct {
	gorm.Model

	Time   time.Time `json:"time"`
	Note   string    `json:"note"`
	UserID uint      `json:"user_id"`

	Activities []Activity `json:"activities" gorm:"foreignKey:DateID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
