package repository

import (
	"github.com/steppbol/activity-manager/model"
)

type UserRepository struct {
	baseRepository *BaseRepository
}

func NewUserRepository(br *BaseRepository) (*UserRepository, error) {
	return &UserRepository{
		baseRepository: br,
	}, nil
}

func (ur UserRepository) Create(user *model.User) {
	ur.baseRepository.database.Create(&user)
}

func (ur UserRepository) Update(user *model.User, update map[string]interface{}) {
	ur.baseRepository.database.Model(&user).Updates(update)
}

func (ur UserRepository) FindByUsername(name string) (*model.User, error) {
	var user model.User

	err := ur.baseRepository.database.Where("username = ?", name).First(&user).Error

	return &user, err
}

func (ur UserRepository) FindByID(id uint) (*model.User, error) {
	var user model.User

	err := ur.baseRepository.database.Where("id = ?", id).First(&user).Error

	return &user, err
}

func (ur UserRepository) DeleteByID(id uint) {
	ur.baseRepository.database.Delete(&model.User{}).Where("id = ?", id)
}
