package model

import "gorm.io/gorm"

type Activity struct {
	gorm.Model

	Title       string `json:"title"`
	Description string `json:"description"`
	Content     string `json:"content"`
	DateID      string `json:"date_id"`

	Tags []Tag `json:"tags" gorm:"many2many:activities_tags;foreignKey:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
