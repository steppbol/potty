package services

import (
	"time"

	"github.com/steppbol/activity-manager/internal/models"
	"github.com/steppbol/activity-manager/internal/repositories"
)

type DateService struct {
	dateRepository *repositories.DateRepository
	userService    *UserService
}

func NewDateService(us *UserService, dr *repositories.DateRepository) (*DateService, error) {
	return &DateService{
		dateRepository: dr,
		userService:    us,
	}, nil
}

func (ds DateService) Create(time time.Time, userId uint, note string) *models.Date {
	_, err := ds.userService.FindByID(userId)
	if err != nil {
		return nil
	}

	fDate, _ := ds.dateRepository.FindAllByTimeAndUserID(userId, time)
	if fDate.ID != 0 {
		return nil
	}

	date := ds.createDate(time, userId, note)
	ds.dateRepository.Create(date)
	return date
}

func (ds DateService) Update(id uint, update map[string]interface{}) *models.Date {
	date, err := ds.FindByID(id)
	if err != nil {
		return nil
	}

	ds.dateRepository.Update(date, update)
	return date
}

func (ds DateService) FindAllByUserID(userId uint) *[]models.Date {
	return ds.dateRepository.FindAllByUserID(userId)
}

func (ds DateService) FindByID(id uint) (*models.Date, error) {
	return ds.dateRepository.FindByID(id)
}

func (ds DateService) DeleteByID(id uint) {
	ds.dateRepository.DeleteByID(id)
}

func (ds DateService) ExportToXLSX() {

}

func (ds DateService) ImportFromXLSX() {
}

func (ds DateService) createDate(time time.Time, userId uint, note string) *models.Date {
	return &models.Date{
		Time:   time,
		Note:   note,
		UserID: userId,
	}
}
