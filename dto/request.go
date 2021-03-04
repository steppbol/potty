package dto

import "time"

type UserDTO struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserUpdateRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type DateDTO struct {
	Time   time.Time `json:"time" binding:"required"`
	Note   string    `json:"note"`
	UserID string    `json:"user_id" binding:"required"`
}

type TagDTO struct {
	Name string `json:"name" binding:"required"`
}
