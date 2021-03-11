package repositories

import (
	"github.com/steppbol/activity-manager/internal/models"
)

type TagRepository struct {
	baseRepository *BaseRepository
}

func NewTagRepository(br *BaseRepository) (*TagRepository, error) {
	return &TagRepository{
		baseRepository: br,
	}, nil
}

func (tr TagRepository) Create(tag *models.Tag) {
	tr.baseRepository.database.Create(&tag)
}

func (tr TagRepository) Update(tag *models.Tag, name string) {
	tr.baseRepository.database.Model(&tag).Update("name", name)
}

func (tr TagRepository) FindAll() *[]models.Tag {
	var tags []models.Tag

	tr.baseRepository.database.Find(&tags)

	return &tags
}

func (tr TagRepository) FindByID(id uint) (*models.Tag, error) {
	var tag models.Tag

	err := tr.baseRepository.database.Where("id = ?", id).First(&tag).Error

	return &tag, err
}

func (tr TagRepository) FindByAllByIDs(ids []uint) *[]models.Tag {
	var tags []models.Tag

	tr.baseRepository.database.Where("id IN (?)", ids).Find(&tags)

	return &tags
}

func (tr TagRepository) FindByName(name string) (*models.Tag, error) {
	var tag models.Tag

	err := tr.baseRepository.database.Where("name = ?", name).First(&tag).Error

	return &tag, err
}

func (tr TagRepository) DeleteByID(id uint) {
	tr.baseRepository.database.Delete(&models.Tag{}).Where("id = ?", id)
}
