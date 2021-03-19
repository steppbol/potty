package models

import "gorm.io/gorm"

type User struct {
	gorm.Model

	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`

	Dates []Date `json:"dates" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
