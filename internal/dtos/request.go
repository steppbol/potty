package dtos

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
	UserID uint      `json:"user_id" binding:"required"`
}

type DateUpdateRequest struct {
	Time time.Time `json:"time" binding:"required"`
	Note string    `json:"note"`
}

type ActivityDTO struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Content     string `json:"content" binding:"required"`
	DateID      uint   `json:"date_id" binding:"required"`
	TagIDs      []uint `json:"tag_ids"`
}

type ActivityUpdateRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Content     string `json:"content"`
	TagIDs      []uint `json:"tag_ids"`
}

type TagDTO struct {
	Name string `json:"name" binding:"required"`
}

type FindByUserIDRequest struct {
	UserID uint `json:"user_id" binding:"required"`
}
