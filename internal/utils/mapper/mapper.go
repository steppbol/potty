package mapper

import (
	"time"

	"github.com/steppbol/activity-manager/internal/dtos"
	"github.com/steppbol/activity-manager/internal/models"
)

func UserUpdateRequestToMap(ur dtos.UserUpdateRequest) *map[string]interface{} {
	m := make(map[string]interface{})

	if ur.Email != "" {
		m["email"] = ur.Email
	}
	if ur.Username != "" {
		m["username"] = ur.Username
	}
	if ur.Password != "" {
		m["password"] = ur.Password
	}

	return &m
}

func DateUpdateRequestToMap(dr dtos.DateUpdateRequest) *map[string]interface{} {
	m := make(map[string]interface{})

	var zTime time.Time
	if dr.Time != zTime {
		m["time"] = dr.Time
	}
	if dr.Note != "" {
		m["note"] = dr.Note
	}

	return &m
}

func ActivityUpdateRequestToMap(ar dtos.ActivityUpdateRequest) *map[string]interface{} {
	m := make(map[string]interface{})

	if ar.Title != "" {
		m["title"] = ar.Title
	}
	if ar.Description != "" {
		m["description"] = ar.Description
	}
	if ar.Content != "" {
		m["content"] = ar.Content
	}
	if len(ar.TagIDs) > 0 {
		m["tag_ids"] = ar.TagIDs
	}
	if ar.Place != "" {
		m["place"] = ar.Place
	}
	if ar.Price != "" {
		m["price"] = ar.Price
	}

	return &m
}

func ActivityDTOToActivity(a dtos.ActivityDTO) *models.Activity {
	return &models.Activity{
		Title:       a.Title,
		Description: a.Description,
		Content:     a.Content,
		Place:       a.Place,
		Price:       a.Price,
	}
}

func ActivityWithDateIDDTOToActivity(a dtos.ActivityWithDateIDDTO) *models.Activity {
	return &models.Activity{
		Title:       a.Title,
		Description: a.Description,
		Content:     a.Content,
		Place:       a.Place,
		Price:       a.Price,
		DateID:      a.DateID,
	}
}
