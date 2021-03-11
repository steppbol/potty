package repositories

import (
	"time"

	"github.com/steppbol/activity-manager/internal/models"
)

type DateRepository struct {
	baseRepository *BaseRepository
}

func NewDateRepository(br *BaseRepository) (*DateRepository, error) {
	return &DateRepository{
		baseRepository: br,
	}, nil
}

func (dr DateRepository) Create(activity *models.Date) {
	dr.baseRepository.database.Create(&activity)
}

func (dr DateRepository) Update(date *models.Date, update map[string]interface{}) {
	dr.baseRepository.database.Model(&date).Updates(update)
}

func (dr DateRepository) FindAllByUserID(userId uint) *[]models.Date {
	var dates []models.Date

	dr.baseRepository.database.Where("user_id = ?", userId).Preload("Activities").Find(&dates)

	return &dates
}

func (dr DateRepository) FindAllByTimeAndUserID(userId uint, time time.Time) (*models.Date, error) {
	var date models.Date

	err := dr.baseRepository.database.Where("user_id = ? AND time = ?", userId, time).Preload("Activities").First(&date).Error

	return &date, err
}

func (dr DateRepository) FindByID(id uint) (*models.Date, error) {
	var date models.Date

	err := dr.baseRepository.database.Where("id = ?", id).Preload("Activities").First(&date).Error

	return &date, err
}

func (dr DateRepository) DeleteByID(id uint) {
	dr.baseRepository.database.Delete(&models.Tag{}).Where("id = ?", id)
}
