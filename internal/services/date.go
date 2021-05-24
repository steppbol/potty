package services

import (
	"io"
	"time"

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

	fDate, _ := ds.FindByTimeAndUserID(userId, time)
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

func (ds DateService) FindAllByUserIDAndNotDeleted(userId uint) *[]models.Date {
	return ds.dateRepository.FindAllByUserIDAndNotDeleted(userId)
}

func (ds DateService) FindByTimeAndUserID(userId uint, time time.Time) (*models.Date, error) {
	return ds.dateRepository.FindByTimeAndUserID(userId, time)
}

func (ds DateService) FindByID(id uint) (*models.Date, error) {
	return ds.dateRepository.FindByID(id)
}

func (ds DateService) DeleteByID(id uint) {
	ds.dateRepository.DeleteByID(id)
}

func (ds DateService) ExportToXLSX(userId uint, from, to *time.Time) (string, error) {
	user, err := ds.userService.FindByID(userId)
	if err != nil {
		return "", nil
	}

	dates := ds.FindAllByUserIDAndNotDeleted(userId)

	return ds.xlsxService.Export(user.Username, *ds.getDatesBetween(*dates, from, to))
}

func (ds DateService) ImportFromXLSX(userId uint, r io.Reader) (*[]models.Date, error) {
	_, err := ds.userService.FindByID(userId)
	if err != nil {
		return nil, err
	}

	return ds.xlsxService.Import(r)
}

func (ds DateService) createDate(time time.Time, userId uint, note string) *models.Date {
	return &models.Date{
		Time:   time,
		Note:   note,
		UserID: userId,
	}
}

func (ds DateService) getDatesBetween(dates []models.Date, from, to *time.Time) *[]models.Date {
	uDates := make([]models.Date, 0)

	if from != nil && to != nil {
		for i := range dates {
			if dates[i].Time.After(*from) && dates[i].Time.Before(*to) {
				uDates = append(uDates, dates[i])
			}
		}
	} else if from != nil {
		for i := range dates {
			if dates[i].Time.After(*from) {
				uDates = append(uDates, dates[i])
			}
		}
	} else if to != nil {
		for i := range dates {
			if dates[i].Time.Before(*to) {
				uDates = append(uDates, dates[i])
			}
		}
	} else {
		uDates = dates
	}

	return &uDates
}
