package dtos

import "time"

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type UserDTO struct {
	Email    string `json:"email" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserUpdateRequest struct {
	Email    string `json:"email"`
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

type ActivityWithDateIDDTO struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Content     string `json:"content" binding:"required"`
	Place       string `json:"place"`
	Price       string `json:"price"`
	DateID      uint   `json:"date_id" binding:"required"`
	TagIDs      []uint `json:"tag_ids"`
}

type ActivityDTO struct {
	Username    string    `json:"username" binding:"required"`
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	Content     string    `json:"content" binding:"required"`
	Place       string    `json:"place"`
	Price       string    `json:"price"`
	Date        time.Time `json:"date" binding:"required"`
	TagIDs      []uint    `json:"tag_ids"`
}

type ActivityUpdateRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Content     string `json:"content"`
	Place       string `json:"place"`
	Price       string `json:"price"`
	TagIDs      []uint `json:"tag_ids"`
}

type TagDTO struct {
	Name string `json:"name" binding:"required"`
}

type UserIDRequest struct {
	UserID uint `json:"user_id" binding:"required"`
}
type ExportToXLSXRequest struct {
	UserID    uint      `json:"user_id" binding:"required"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

type FindByTagsRequest struct {
	UserID uint   `json:"user_id" binding:"required"`
	TagIDs []uint `json:"tag_ids" binding:"required"`
}
