package mapper

import (
	"github.com/steppbol/activity-manager/dto"
)

func UserUpdateRequestToMap(ur dto.UserUpdateRequest) *map[string]interface{} {
	m := make(map[string]interface{})

	if ur.Username != "" {
		m["username"] = ur.Username
	}
	if ur.Password != "" {
		m["password"] = ur.Password
	}

	return &m
}
