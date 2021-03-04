package model

import (
	"time"

	"gorm.io/gorm"
)

type Date struct {
	gorm.Model

	Time   time.Time
	Note   string `json:"note"`
	UserID string `json:"user_id"`

	Activities []Activity `json:"activities" gorm:"foreignKey:DateID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
