package repositories

import "github.com/steppbol/activity-manager/internal/models"

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

func (ar ActivityRepository) Update(activity *models.Activity) {
	ar.baseRepository.database.Save(activity)
}

func (ar ActivityRepository) FindAllByUserID(userId uint) *[]models.Activity {
	var activities []models.Activity

	ar.baseRepository.database.Where("user_id = ?", userId).Preload("Tags").Find(&activities)

	return &activities
}

func (ar ActivityRepository) FindAllByTagsAndDateID(dateId uint, tagIds []uint) *[]models.Activity {
	var activities []models.Activity

	ar.baseRepository.database.Preload("tag_id IN (?)", tagIds).Where("date_id = ?", dateId).Find(&activities)

	return &activities
}

func (ar ActivityRepository) FindByID(id uint) (*models.Activity, error) {
	var activity models.Activity

	err := ar.baseRepository.database.Where("id = ?", id).Preload("Tags").First(&activity).Error

	return &activity, err
}

func (ar ActivityRepository) DeleteByID(id uint) {
	ar.baseRepository.database.Delete(&models.Activity{}).Where("id = ?", id)
}
