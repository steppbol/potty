package repositories

import (
	"github.com/steppbol/activity-manager/internal/models"
)

type UserRepository struct {
	baseRepository *BaseRepository
}

func NewUserRepository(br *BaseRepository) *UserRepository {
	return &UserRepository{
		baseRepository: br,
	}
}

func (ur UserRepository) Create(user *models.User) {
	ur.baseRepository.database.Create(&user)
}

func (ur UserRepository) Update(user *models.User, update map[string]interface{}) {
	ur.baseRepository.database.Model(&user).Updates(update)
}

func (ur UserRepository) FindByUsername(name string) (*models.User, error) {
	var user models.User

	err := ur.baseRepository.database.Where("username = ?", name).First(&user).Error

	return &user, err
}

func (ur UserRepository) FindByID(id uint) (*models.User, error) {
	var user models.User

	err := ur.baseRepository.database.Where("id = ?", id).First(&user).Error

	return &user, err
}

func (ur UserRepository) DeleteByID(id uint) {
	ur.baseRepository.database.Delete(&models.User{}).Where("id = ?", id)
}
