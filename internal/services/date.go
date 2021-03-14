package services

import (
	"time"

	"github.com/tealeg/xlsx"

	"github.com/steppbol/activity-manager/internal/models"
	"github.com/steppbol/activity-manager/internal/repositories"
)

type DateService struct {
	dateRepository *repositories.DateRepository
	userService    *UserService
	xlsxService    *XLSXService
}

func NewDateService(us *UserService, xs *XLSXService, dr *repositories.DateRepository) *DateService {
	return &DateService{
		dateRepository: dr,
		userService:    us,
		xlsxService:    xs,
	}
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

func (ds DateService) ExportToXLSX(userId uint) (*xlsx.File, error) {
	user, err := ds.userService.FindByID(userId)
	if err != nil {
		return nil, nil
	}

	dates := ds.FindAllByUserID(userId)

	activities := make([]models.Activity, 0)

	for i := range *dates {
		activities = append(activities, (*dates)[i].Activities...)
	}

	tags := make(map[uint][]models.Tag)

	for i := range activities {
		tags[activities[i].ID] = activities[i].Tags
	}

	return ds.xlsxService.Export(user.Username, *dates, activities, tags)
}

func (ds DateService) createDate(time time.Time, userId uint, note string) *models.Date {
	return &models.Date{
		Time:   time,
		Note:   note,
		UserID: userId,
	}
}
