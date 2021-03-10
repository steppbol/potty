package services

import (
	"github.com/steppbol/activity-manager/internal/models"
	"github.com/steppbol/activity-manager/internal/repositories"
)

type TagService struct {
	tagRepository *repositories.TagRepository
}

func NewTagService(tr *repositories.TagRepository) (*TagService, error) {
	return &TagService{
		tagRepository: tr,
	}, nil
}

func (ts TagService) Create(name string) *models.Tag {
	fTag, _ := ts.tagRepository.FindByName(name)
	if fTag.ID != 0 {
		return fTag
	}

	tag := ts.createTag(name)
	ts.tagRepository.Create(tag)
	return tag
}

func (ts TagService) Update(id uint, name string) *models.Tag {
	tag, err := ts.FindByID(id)
	if err != nil {
		return nil
	}

	ts.tagRepository.Update(tag, name)
	return tag
}

func (ts TagService) FindAll() *[]models.Tag {
	return ts.tagRepository.FindAll()
}

func (ts TagService) FindByID(id uint) (*models.Tag, error) {
	return ts.tagRepository.FindByID(id)
}

func (ts TagService) DeleteByID(id uint) {
	ts.tagRepository.DeleteByID(id)
}

func (ts TagService) createTag(name string) *models.Tag {
	return &models.Tag{
		Name: name,
	}
}
