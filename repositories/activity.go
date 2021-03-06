package repositories

import "github.com/steppbol/activity-manager/models"

type ActivityRepository struct {
	baseRepository *BaseRepository
}

func NewActivityRepository(br *BaseRepository) (*ActivityRepository, error) {
	return &ActivityRepository{
		baseRepository: br,
	}, nil
}

func (ar ActivityRepository) Create(activity *models.Activity) {
	ar.baseRepository.database.Create(&activity)
}

func (ar ActivityRepository) Update(id uint, activity *models.Activity) {
	ar.baseRepository.database.Model(&models.Activity{}).Where("id = ?", id).Updates(activity)
}

func (ar ActivityRepository) FindAll() *[]models.Activity {
	var activities []models.Activity

	ar.baseRepository.database.Find(&activities)

	return &activities
}

func (ar ActivityRepository) FindById(id uint) (*models.Activity, error) {
	var activity models.Activity

	err := ar.baseRepository.database.Where("id = ?", id).First(&activity).Error

	return &activity, err
}

func (ar ActivityRepository) FindByTags(id uint, tags []models.Activity) *[]models.Activity {
	var activities []models.Activity

	ar.baseRepository.database.Where("id = ?", id).Where("tag IN ?", tags).Find(&activities)

	return &activities
}

func (ar ActivityRepository) DeleteById(id uint) {
	ar.baseRepository.database.Delete(&models.Activity{}).Where("id = ?", id)
}
