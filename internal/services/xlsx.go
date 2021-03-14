package services

import (
	"fmt"
	"strconv"
	"time"

	"github.com/steppbol/activity-manager/configs"
	"github.com/steppbol/activity-manager/internal/models"

	"github.com/tealeg/xlsx"
)

const xlsxExtension = ".xlsx"

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

func (xs XLSXService) Export(username string, dates []models.Date, activities []models.Activity, tags map[uint][]models.Tag) (*xlsx.File, error) {
	file := xlsx.NewFile()
	dSheet, err := file.AddSheet("Dates")
	if err != nil {
		return nil, err
	}

	xs.exportDates(dSheet, dates)

	aSheet, err := file.AddSheet("Activities")
	if err != nil {
		return nil, err
	}

	xs.exportActivities(aSheet, activities)

	tSheet, err := file.AddSheet("Tags")
	if err != nil {
		return nil, err
	}

	xs.exportTags(tSheet, tags)

	currTime := strconv.Itoa(int(time.Now().Unix()))
	filename := username + currTime + xlsxExtension

	err = file.Save(xs.config.XLSXExportPath + filename)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (xs XLSXService) exportDates(sheet *xlsx.Sheet, dates []models.Date) {
	titles := []string{"ID", "CreatedAt", "UpdatedAt", "DeletedAt", "Time", "Note", "UserID"}

	row := sheet.AddRow()
	var cell *xlsx.Cell
	for i := range titles {
		cell = row.AddCell()
		cell.Value = titles[i]
	}

	for i := range dates {
		values := []string{
			fmt.Sprintf("%d", dates[i].ID),
			dates[i].CreatedAt.String(),
			dates[i].UpdatedAt.String(),
			dates[i].DeletedAt.Time.String(),
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
	titles := []string{"ID", "CreatedAt", "UpdatedAt", "DeletedAt", "Title", "Description", "Content", "DateID"}

	row := sheet.AddRow()
	var cell *xlsx.Cell
	for i := range titles {
		cell = row.AddCell()
		cell.Value = titles[i]
	}

	for i := range activities {
		values := []string{
			fmt.Sprintf("%d", activities[i].ID),
			activities[i].CreatedAt.String(),
			activities[i].UpdatedAt.String(),
			activities[i].DeletedAt.Time.String(),
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
	titles := []string{"ID", "ActivityID", "CreatedAt", "UpdatedAt", "DeletedAt", "Name"}

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
				v[i].CreatedAt.String(),
				v[i].UpdatedAt.String(),
				v[i].DeletedAt.Time.String(),
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
