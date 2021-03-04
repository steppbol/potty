package service

import (
	"github.com/steppbol/activity-manager/model"
	"github.com/steppbol/activity-manager/repository"
)

type UserService struct {
	userRepository *repository.UserRepository
}

func NewUserService(ur *repository.UserRepository) (*UserService, error) {
	return &UserService{
		userRepository: ur,
	}, nil
}

func (us UserService) Create(username, password string) *model.User {
	fUser, _ := us.userRepository.FindByUsername(username)
	if fUser.Username != "" || fUser.Password != "" {
		return nil
	}

	user := us.createUser(username, password)
	us.userRepository.Create(user)
	return user
}

func (us UserService) Update(id uint, update map[string]interface{}) *model.User {
	user, err := us.FindByID(id)
	if err != nil {
		return nil
	}

	us.userRepository.Update(user, update)
	return user
}

func (us UserService) FindByID(id uint) (*model.User, error) {
	return us.userRepository.FindByID(id)
}

func (us UserService) DeleteByID(id uint) {
	us.userRepository.DeleteByID(id)
}

func (us UserService) createUser(username, password string) *model.User {
	return &model.User{
		Username: username,
		Password: password,
	}
}
