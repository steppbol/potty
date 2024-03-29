package api

import (
	"io"
	"os"
	"time"

	"github.com/steppbol/activity-manager/internal/services"
)

type XLSXBaseAPI struct {
	TagService      *services.TagService
	ActivityService *services.ActivityService
	DateService     *services.DateService
}

func NewXLSXBaseAPI(ts *services.TagService, as *services.ActivityService, ds *services.DateService) *XLSXBaseAPI {
	return &XLSXBaseAPI{
		TagService:      ts,
		ActivityService: as,
		DateService:     ds,
	}
}

func (ba XLSXBaseAPI) ExportToXLSX(userId uint, from, to time.Time) (string, error) {
	return ba.DateService.ExportToXLSX(userId, &from, &to)
}

func (ba XLSXBaseAPI) ImportFromXLSX(userId uint, r io.Reader) error {
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

			activity.DateID = cDate.ID
			ba.ActivityService.CreateWithDateID(activity, cTags)
		}
	}

	return nil
}

func (ba XLSXBaseAPI) DeleteStaticData(path string) error {
	return os.Remove(path)
}
