package service

import (
	"github.com/steppbol/activity-manager/model"
	"github.com/steppbol/activity-manager/repository"
)

type TagService struct {
	tagRepository *repository.TagRepository
}

func NewTagService(tr *repository.TagRepository) (*TagService, error) {
	return &TagService{
		tagRepository: tr,
	}, nil
}

func (ts TagService) Create(name string) *model.Tag {
	fTag, err := ts.tagRepository.FindByName(name)
	if err != nil {
		return nil
	}

	if fTag != nil {
		return fTag
	}

	tag := ts.createTag(name)
	ts.tagRepository.Create(tag)
	return tag
}

func (ts TagService) Update(id uint, name string) *model.Tag {
	tag, err := ts.FindByID(id)
	if err != nil {
		return nil
	}

	ts.tagRepository.Update(tag, name)
	return tag
}

func (ts TagService) FindAll() *[]model.Tag {
	return ts.tagRepository.FindAll()
}

func (ts TagService) FindByID(id uint) (*model.Tag, error) {
	return ts.tagRepository.FindByID(id)
}

func (ts TagService) DeleteByID(id uint) {
	ts.tagRepository.DeleteByID(id)
}

func (ts TagService) createTag(name string) *model.Tag {
	return &model.Tag{
		Name: name,
	}
}
