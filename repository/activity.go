package repository

import "github.com/steppbol/activity-manager/model"

type ActivityRepository struct {
	baseRepository *BaseRepository
}

func NewActivityRepository(br *BaseRepository) (*ActivityRepository, error) {
	return &ActivityRepository{
		baseRepository: br,
	}, nil
}

func (ar ActivityRepository) Create(activity *model.Activity) {
	ar.baseRepository.database.Create(&activity)
}

func (ar ActivityRepository) Update(id uint, activity *model.Activity) {
	ar.baseRepository.database.Model(&model.Activity{}).Where("id = ?", id).Updates(activity)
}

func (ar ActivityRepository) FindAll() *[]model.Activity {
	var activities []model.Activity

	ar.baseRepository.database.Find(&activities)

	return &activities
}

func (ar ActivityRepository) FindById(id uint) (*model.Activity, error) {
	var activity model.Activity

	err := ar.baseRepository.database.Where("id = ?", id).First(&activity).Error

	return &activity, err
}

func (ar ActivityRepository) FindByTags(id uint, tags []model.Activity) *[]model.Activity {
	var activities []model.Activity

	ar.baseRepository.database.Where("id = ?", id).Where("tag IN ?", tags).Find(&activities)

	return &activities
}

func (ar ActivityRepository) DeleteById(id uint) {
	ar.baseRepository.database.Delete(&model.Activity{}).Where("id = ?", id)
}
