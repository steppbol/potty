package services

import (
	"time"

	"github.com/steppbol/activity-manager/internal/models"
	"github.com/steppbol/activity-manager/internal/repositories"
)

type ActivityService struct {
	activityRepository *repositories.ActivityRepository
	tagService         *TagService
	dateService        *DateService
	userService        *UserService
}

func NewActivityService(ts *TagService, ds *DateService, us *UserService, ar *repositories.ActivityRepository) *ActivityService {
	return &ActivityService{
		activityRepository: ar,
		tagService:         ts,
		dateService:        ds,
		userService:        us,
	}
}

func (as ActivityService) Create(activity models.Activity, username string, date time.Time, tagIds []uint) *models.Activity {
	user, _ := as.userService.FindByUsername(username)
	if user.ID == 0 {
		return nil
	}

	fDate, _ := as.dateService.FindByTimeAndUserID(user.ID, date)

	cDate := fDate
	if fDate.ID == 0 {
		cDate = as.dateService.Create(date, user.ID, "")
	}

	tags := as.tagService.FindAllByIDs(tagIds)
	cActivity := as.createActivity(activity, cDate.ID, *tags)

	as.activityRepository.Create(cActivity)
	return cActivity
}

func (as ActivityService) CreateWithDateID(activity models.Activity, tagIds []uint) *models.Activity {
	_, err := as.dateService.FindByID(activity.DateID)
	if err != nil {
		return nil
	}

	tags := as.tagService.FindAllByIDs(tagIds)
	cActivity := as.createActivity(activity, activity.DateID, *tags)

	as.activityRepository.Create(cActivity)
	return cActivity
}

func (as ActivityService) Update(id uint, update map[string]interface{}) *models.Activity {
	activity, err := as.activityRepository.FindByID(id)
	if err != nil {
		return nil
	}

	if update["title"].(string) != "" {
		activity.Title = update["title"].(string)
	}
	if update["description"].(string) != "" {
		activity.Description = update["description"].(string)
	}
	if update["content"].(string) != "" {
		activity.Content = update["content"].(string)
	}
	if len(update["tag_ids"].([]uint)) > 0 {
		activity.Tags = *as.tagService.FindAllByIDs(update["tag_ids"].([]uint))
	}
	if update["place"].(string) != "" {
		activity.Place = update["place"].(string)
	}
	if update["price"].(string) != "" {
		activity.Price = update["price"].(string)
	}

	as.activityRepository.Update(activity)
	return activity
}

func (as ActivityService) FindAllByUserID(userId uint) *[]models.Activity {
	return as.activityRepository.FindAllByUserID(userId)
}

func (as ActivityService) FindAllByTags(userId uint, tagIds []uint) *[]models.Activity {
	return as.activityRepository.FindAllByTagsAndUserID(userId, tagIds)
}

func (as ActivityService) DeleteByID(id uint) {
	as.activityRepository.DeleteByID(id)
}

func (as ActivityService) createActivity(activity models.Activity, dateId uint, tags []models.Tag) *models.Activity {
	return &models.Activity{
		Title:       activity.Title,
		Description: activity.Description,
		Content:     activity.Content,
		Place:       activity.Place,
		Price:       activity.Price,
		DateID:      dateId,
		Tags:        tags,
	}
}
