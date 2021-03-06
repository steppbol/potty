package mapper

import (
	"time"

	"github.com/steppbol/activity-manager/dtos"
)

func UserUpdateRequestToMap(ur dtos.UserUpdateRequest) *map[string]interface{} {
	m := make(map[string]interface{})

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
