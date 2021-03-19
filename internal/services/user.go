package services

import (
	"golang.org/x/crypto/bcrypt"

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

func (us UserService) Create(username, password, email string) *models.User {
	fUser, _ := us.FindByUsername(username)
	if fUser.ID != 0 {
		return nil
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := us.createUser(username, string(hashedPassword), email)
	us.userRepository.Create(user)
	return user
}

func (us UserService) Update(id uint, update map[string]interface{}) *models.User {
	user, err := us.FindByID(id)
	if err != nil {
		return nil
	}

	if update["password"] != "" {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(update["password"].(string)), bcrypt.DefaultCost)
		update["password"] = string(hashedPassword)
	}

	us.userRepository.Update(user, update)
	return user
}

func (us UserService) FindByID(id uint) (*models.User, error) {
	return us.userRepository.FindByID(id)
}

func (us UserService) FindByUsername(username string) (*models.User, error) {
	return us.userRepository.FindByUsername(username)
}

func (us UserService) DeleteByID(id uint) {
	us.userRepository.DeleteByID(id)
}

func (us UserService) CheckUser(username, password string) bool {
	user, _ := us.FindByUsername(username)
	if user.ID == 0 {
		return false
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return false
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		return false
	}

	return true
}

func (us UserService) createUser(username, password, email string) *models.User {
	return &models.User{
		Username: username,
		Password: password,
		Email:    email,
	}
}
