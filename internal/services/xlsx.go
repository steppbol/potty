package services

import (
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"gorm.io/gorm"

	"github.com/steppbol/activity-manager/configs"
	"github.com/steppbol/activity-manager/internal/models"

	"github.com/tealeg/xlsx"
)

const (
	xlsxExtension = ".xlsx"
	timeLayout    = "2006-01-02 15:04:05.999999 -0700 +03"
)

type XLSXService struct {
	exportPath string
	config     configs.Application
}

func NewXLSXService(conf *configs.Application) *XLSXService {
	return &XLSXService{
		exportPath: conf.XLSXExportPath,
		config:     *conf,
	}
}

func (xs XLSXService) Export(username string, dates []models.Date) (string, error) {
	activities := make([]models.Activity, 0)

	for i := range dates {
		activities = append(activities, dates[i].Activities...)
	}

	tags := make(map[uint][]models.Tag)

	for i := range activities {
		tags[activities[i].ID] = activities[i].Tags
	}

	file := xlsx.NewFile()
	dSheet, err := file.AddSheet("Dates")
	if err != nil {
		return "", err
	}

	xs.exportDates(dSheet, dates)

	aSheet, err := file.AddSheet("Activities")
	if err != nil {
		return "", err
	}

	xs.exportActivities(aSheet, activities)

	tSheet, err := file.AddSheet("Tags")
	if err != nil {
		return "", err
	}

	xs.exportTags(tSheet, tags)

	currTime := strconv.Itoa(int(time.Now().Unix()))
	filename := username + currTime + xlsxExtension
	path := xs.config.XLSXExportPath + filename

	err = file.Save(path)
	if err != nil {
		return "", err
	}

	return path, nil
}

func (xs XLSXService) Import(r io.Reader) (*[]models.Date, error) {
	file, err := excelize.OpenReader(r)
	if err != nil {
		return nil, err
	}

	rows, err := file.GetRows("Tags")
	if err != nil {
		return nil, err
	}

	tags, err := xs.createTags(&rows)
	if err != nil {
		return nil, err
	}

	rows, err = file.GetRows("Activities")
	if err != nil {
		return nil, err
	}

	activities, err := xs.createActivities(&rows, tags)
	if err != nil {
		return nil, err
	}

	rows, err = file.GetRows("Dates")
	if err != nil {
		return nil, err
	}

	dates, err := xs.createDates(&rows, activities)
	if err != nil {
		return nil, err
	}

	return dates, nil
}

func (xs XLSXService) exportDates(sheet *xlsx.Sheet, dates []models.Date) {
	titles := []string{"ID", "Time", "Note", "UserID"}

	row := sheet.AddRow()
	var cell *xlsx.Cell
	for i := range titles {
		cell = row.AddCell()
		cell.Value = titles[i]
	}

	for i := range dates {
		values := []string{
			fmt.Sprintf("%d", dates[i].ID),
			dates[i].Time.String(),
			dates[i].Note,
			fmt.Sprintf("%d", dates[i].UserID),
		}

		row = sheet.AddRow()
		for vi := range values {
			cell = row.AddCell()
			cell.Value = values[vi]
		}
	}
}

func (xs XLSXService) exportActivities(sheet *xlsx.Sheet, activities []models.Activity) {
	titles := []string{"ID", "Title", "Description", "Content", "DateID"}

	row := sheet.AddRow()
	var cell *xlsx.Cell
	for i := range titles {
		cell = row.AddCell()
		cell.Value = titles[i]
	}

	for i := range activities {
		values := []string{
			fmt.Sprintf("%d", activities[i].ID),
			activities[i].Title,
			activities[i].Description,
			activities[i].Content,
			fmt.Sprintf("%d", activities[i].DateID),
		}

		row = sheet.AddRow()
		for vi := range values {
			cell = row.AddCell()
			cell.Value = values[vi]
		}
	}
}

func (xs XLSXService) exportTags(sheet *xlsx.Sheet, tags map[uint][]models.Tag) {
	titles := []string{"ID", "ActivityID", "Name"}

	row := sheet.AddRow()
	var cell *xlsx.Cell
	for i := range titles {
		cell = row.AddCell()
		cell.Value = titles[i]
	}

	for k, v := range tags {
		for i := range v {
			values := []string{
				fmt.Sprintf("%d", v[i].ID),
				fmt.Sprintf("%d", k),
				v[i].Name,
			}

			row = sheet.AddRow()
			for vi := range values {
				cell = row.AddCell()
				cell.Value = values[vi]
			}
		}
	}
}

func (xs XLSXService) createTags(rows *[][]string) (*map[uint][]models.Tag, error) {
	activityTags := make(map[uint][]models.Tag)

	for i, row := range *rows {
		if i > 0 {
			var data []string
			for ir := range row {
				data = append(data, row[ir])
			}

			if len(data) > 0 {
				activityId, err := strconv.Atoi(data[1])
				if err != nil {
					return nil, err
				}

				tag, err := xs.createTag(data)
				if err != nil {
					return nil, err
				}

				tags := activityTags[uint(activityId)]

				if tags == nil {
					tags = make([]models.Tag, 0)
				}

				tags = append(tags, *tag)

				activityTags[uint(activityId)] = tags
			}
		}
	}

	return &activityTags, nil
}

func (xs XLSXService) createActivities(rows *[][]string, tags *map[uint][]models.Tag) (*map[uint][]models.Activity, error) {
	dateActivities := make(map[uint][]models.Activity)

	for i, row := range *rows {
		if i > 0 {
			var data []string
			for ir := range row {
				data = append(data, row[ir])
			}

			if len(data) > 0 {
				dateId, err := strconv.Atoi(data[4])
				if err != nil {
					return nil, err
				}

				activityId, err := strconv.Atoi(data[0])
				if err != nil {
					return nil, err
				}

				aTags := (*tags)[uint(activityId)]

				activity, err := xs.createActivity(data, &aTags)
				if err != nil {
					return nil, err
				}

				activities := dateActivities[uint(dateId)]

				if activities == nil {
					activities = make([]models.Activity, 0)
				}

				activities = append(activities, *activity)

				dateActivities[uint(dateId)] = activities
			}
		}
	}

	return &dateActivities, nil
}

func (xs XLSXService) createDates(rows *[][]string, activities *map[uint][]models.Activity) (*[]models.Date, error) {
	dates := make([]models.Date, 0)

	for i, row := range *rows {
		if i > 0 {
			var data []string
			for ir := range row {
				data = append(data, row[ir])
			}

			if len(data) > 0 {
				dateId, err := strconv.Atoi(data[0])
				if err != nil {
					return nil, err
				}

				dActivities := (*activities)[uint(dateId)]

				date, err := xs.createDate(data, &dActivities)
				if err != nil {
					return nil, err
				}

				dates = append(dates, *date)
			}
		}
	}

	return &dates, nil
}

func (xs XLSXService) createTag(data []string) (*models.Tag, error) {
	cId, err := strconv.Atoi(data[0])
	if err != nil {
		return nil, err
	}

	return &models.Tag{
		Model: gorm.Model{
			ID: uint(cId),
		},
		Name: data[2],
	}, nil
}

func (xs XLSXService) createActivity(data []string, tags *[]models.Tag) (*models.Activity, error) {
	cId, err := strconv.Atoi(data[0])
	if err != nil {
		return nil, err
	}

	dateId, err := strconv.Atoi(data[4])
	if err != nil {
		return nil, err
	}

	return &models.Activity{
		Model: gorm.Model{
			ID: uint(cId),
		},
		Title:       data[1],
		Description: data[2],
		Content:     data[3],
		DateID:      uint(dateId),
		Tags:        *tags,
	}, nil
}

func (xs XLSXService) createDate(data []string, activities *[]models.Activity) (*models.Date, error) {
	cId, err := strconv.Atoi(data[0])
	if err != nil {
		return nil, err
	}

	cTime, err := time.Parse(timeLayout, data[1])
	if err != nil {
		return nil, err
	}

	userId, err := strconv.Atoi(data[3])
	if err != nil {
		return nil, err
	}

	return &models.Date{
		Model: gorm.Model{
			ID: uint(cId),
		},
		Time:       cTime,
		Note:       data[2],
		UserID:     uint(userId),
		Activities: *activities,
	}, nil
}
