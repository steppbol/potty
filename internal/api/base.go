package api

import (
	"io"
	"os"

	"github.com/steppbol/activity-manager/internal/services"
)

type BaseAPI struct {
	TagService      *services.TagService
	ActivityService *services.ActivityService
	DateService     *services.DateService
	UserService     *services.UserService
}

func NewBaseAPI(ts *services.TagService, as *services.ActivityService, ds *services.DateService, us *services.UserService) *BaseAPI {
	return &BaseAPI{
		TagService:      ts,
		ActivityService: as,
		DateService:     ds,
		UserService:     us,
	}
}

func (ba BaseAPI) ExportToXLSX(userId uint) (string, error) {
	return ba.DateService.ExportToXLSX(userId)
}

func (ba BaseAPI) ImportFromXLSX(userId uint, r io.Reader) error {
	dates, err := ba.DateService.ImportFromXLSX(userId, r)
	if err != nil {
		return err
	}

	for _, date := range *dates {
		cDate := ba.DateService.Create(date.Time, date.UserID, date.Note)
		if cDate == nil {
			return nil
		}

		for _, activity := range date.Activities {
			cTags := make([]uint, 0)
			for _, tag := range activity.Tags {
				ba.TagService.Create(tag.Name)
				cTags = append(cTags, tag.ID)
			}

			ba.ActivityService.Create(activity.Title, activity.Description, activity.Content, cDate.ID, cTags)
		}
	}

	return nil
}

func (ba BaseAPI) DeleteStaticData(path string) error {
	return os.Remove(path)
}
