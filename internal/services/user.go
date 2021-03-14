package services

import (
	"github.com/steppbol/activity-manager/internal/models"
	"github.com/steppbol/activity-manager/internal/repositories"
)

type UserService struct {
	userRepository *repositories.UserRepository
}

func NewUserService(ur *repositories.UserRepository) *UserService {
	return &UserService{
		userRepository: ur,
	}
}

func (us UserService) Create(username, password string) *models.User {
	fUser, _ := us.userRepository.FindByUsername(username)
	if fUser.ID != 0 {
		return nil
	}

	user := us.createUser(username, password)
	us.userRepository.Create(user)
	return user
}

func (us UserService) Update(id uint, update map[string]interface{}) *models.User {
	user, err := us.FindByID(id)
	if err != nil {
		return nil
	}

	us.userRepository.Update(user, update)
	return user
}

func (us UserService) FindByID(id uint) (*models.User, error) {
	return us.userRepository.FindByID(id)
}

func (us UserService) DeleteByID(id uint) {
	us.userRepository.DeleteByID(id)
}

func (us UserService) createUser(username, password string) *models.User {
	return &models.User{
		Username: username,
		Password: password,
	}
}
