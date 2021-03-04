package service

import (
	"time"

	"github.com/steppbol/activity-manager/model"
	"github.com/steppbol/activity-manager/repository"
)

type DateService struct {
	dateRepository *repository.DateRepository
}

func NewDateService(dr *repository.DateRepository) (*DateService, error) {
	return &DateService{
		dateRepository: dr,
	}, nil
}

func (ds DateService) Create(time time.Time, userId, note string) *model.Date {
	fDate, err := ds.dateRepository.FindByTimeAndUserID(userId, time)
	if err != nil {
		return nil
	}

	if fDate != nil {
		return fDate
	}

	date := ds.createDate(time, userId, note)
	ds.dateRepository.Create(date)
	return date
}

func (ds DateService) Update(id uint, update map[string]interface{}) *model.Date {
	date, err := ds.FindById(id)
	if err != nil {
		return nil
	}

	ds.dateRepository.Update(date, update)
	return date
}

func (ds DateService) FindAllByUserID(userId string) *[]model.Date {
	return ds.dateRepository.FindAllByUserID(userId)
}

func (ds DateService) FindById(id uint) (*model.Date, error) {
	return ds.dateRepository.FindByID(id)
}

func (ds DateService) FindAllActivities(id uint) *[]model.Activity {
	return ds.dateRepository.FindAllActivities(id)
}

func (ds DateService) DeleteByID(id uint) {
	ds.dateRepository.DeleteByID(id)
}

func (ds DateService) createDate(time time.Time, userId, note string) *model.Date {
	return &model.Date{
		Time:   time,
		Note:   note,
		UserID: userId,
	}
}
