package repositories

import (
	"gorm.io/gorm/clause"

	"github.com/steppbol/activity-manager/internal/models"
)

type ActivityRepository struct {
	baseRepository *BaseRepository
}

func NewActivityRepository(br *BaseRepository) *ActivityRepository {
	return &ActivityRepository{
		baseRepository: br,
	}
}

func (ar ActivityRepository) Create(activity *models.Activity) {
	ar.baseRepository.database.Create(&activity)
}

func (ar ActivityRepository) Update(activity *models.Activity) {
	ar.baseRepository.database.Save(activity)
}

func (ar ActivityRepository) FindAllByUserID(userId uint) *[]models.Activity {
	var activities []models.Activity

	ar.baseRepository.database.Where("user_id = ?", userId).Preload(clause.Associations).Find(&activities)

	return &activities
}

func (ar ActivityRepository) FindAllByTagsAndUserID(userId uint, tagIds []uint) *[]models.Activity {
	var activities []models.Activity

	ar.baseRepository.database.Joins("JOIN dates ON dates.id = activities.date_id").Joins("JOIN activities_tags ON activities_tags.activity_id=activities.id").Where("user_id = ? AND tag_id IN (?)", userId, tagIds).Find(&activities)

	return &activities
}

func (ar ActivityRepository) FindByID(id uint) (*models.Activity, error) {
	var activity models.Activity

	err := ar.baseRepository.database.Where("id = ?", id).Preload(clause.Associations).First(&activity).Error

	return &activity, err
}

func (ar ActivityRepository) DeleteByID(id uint) {
	ar.baseRepository.database.Delete(&models.Activity{}).Where("id = ?", id)
}
