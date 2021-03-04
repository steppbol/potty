package repository

import (
	"time"

	"github.com/steppbol/activity-manager/model"
)

type DateRepository struct {
	baseRepository *BaseRepository
}

func NewDateRepository(br *BaseRepository) (*DateRepository, error) {
	return &DateRepository{
		baseRepository: br,
	}, nil
}

func (dr DateRepository) Create(activity *model.Date) {
	dr.baseRepository.database.Create(&activity)
}

func (dr DateRepository) Update(date *model.Date, update map[string]interface{}) {
	dr.baseRepository.database.Model(&date).Updates(update)
}

func (dr DateRepository) FindAllByUserID(userId string) *[]model.Date {
	var dates []model.Date

	dr.baseRepository.database.Where("user_id = ?", userId).Find(&dates)

	return &dates
}

func (dr DateRepository) FindByTimeAndUserID(userId string, time time.Time) (*model.Date, error) {
	var date model.Date

	err := dr.baseRepository.database.Where("user_id = ? AND time = ?", userId, time).First(&date).Error

	return &date, err
}

func (dr DateRepository) FindByID(id uint) (*model.Date, error) {
	var date model.Date

	err := dr.baseRepository.database.Where("id = ?", id).First(&date).Error

	return &date, err
}

func (dr DateRepository) FindAllActivities(id uint) *[]model.Activity {
	var activities []model.Activity

	dr.baseRepository.database.Where("id = ?", id).Find(&activities)

	return &activities
}

func (dr DateRepository) DeleteByID(id uint) {
	dr.baseRepository.database.Delete(&model.Tag{}).Where("id = ?", id)
}
