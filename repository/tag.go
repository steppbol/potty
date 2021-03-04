package repository

import (
	"github.com/steppbol/activity-manager/model"
)

type TagRepository struct {
	baseRepository *BaseRepository
}

func NewTagRepository(br *BaseRepository) (*TagRepository, error) {
	return &TagRepository{
		baseRepository: br,
	}, nil
}

func (tr TagRepository) Create(tag *model.Tag) {
	tr.baseRepository.database.Create(&tag)
}

func (tr TagRepository) Update(tag *model.Tag, name string) {
	tr.baseRepository.database.Model(&tag).Update("name", name)
}

func (tr TagRepository) FindAll() *[]model.Tag {
	var tags []model.Tag

	tr.baseRepository.database.Find(&tags)

	return &tags
}

func (tr TagRepository) FindByID(id uint) (*model.Tag, error) {
	var tag model.Tag

	err := tr.baseRepository.database.Where("id = ?", id).First(&tag).Error

	return &tag, err
}

func (tr TagRepository) FindByName(name string) (*model.Tag, error) {
	var tag model.Tag

	err := tr.baseRepository.database.Where("name = ?", name).First(&tag).Error

	return &tag, err
}

func (tr TagRepository) DeleteByID(id uint) {
	tr.baseRepository.database.Delete(&model.Tag{}).Where("id = ?", id)
}
