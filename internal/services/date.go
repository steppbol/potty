package services

import (
	"io"
	"os"
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

func (ds DateService) FindAllByUserIDAndNotDeleted(userId uint) *[]models.Date {
	return ds.dateRepository.FindAllByUserIDAndNotDeleted(userId)
}

func (ds DateService) FindByID(id uint) (*models.Date, error) {
	return ds.dateRepository.FindByID(id)
}

func (ds DateService) DeleteByID(id uint) {
	ds.dateRepository.DeleteByID(id)
}

func (ds DateService) ExportToXLSX(userId uint) (string, error) {
	user, err := ds.userService.FindByID(userId)
	if err != nil {
		return "", nil
	}

	dates := ds.FindAllByUserIDAndNotDeleted(userId)

	return ds.xlsxService.Export(user.Username, *dates)
}

func (ds DateService) ImportFromXLSX(userId uint, r io.Reader) (*[]models.Date, error) {
	_, err := ds.userService.FindByID(userId)
	if err != nil {
		return nil, nil
	}

	return ds.xlsxService.Import(r)
}

func (ds DateService) DeleteStaticData(path string) error {
	return os.Remove(path)
}

func (ds DateService) createDate(time time.Time, userId uint, note string) *models.Date {
	return &models.Date{
		Time:   time,
		Note:   note,
		UserID: userId,
	}
}
