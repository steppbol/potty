package model

import "gorm.io/gorm"

type User struct {
	gorm.Model

	Username string `json:"username"`
	Password string `json:"password"`

	Dates []Date `json:"dates" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
