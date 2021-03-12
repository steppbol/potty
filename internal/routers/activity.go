package routers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/steppbol/activity-manager/internal/dtos"
	"github.com/steppbol/activity-manager/internal/services"
	"github.com/steppbol/activity-manager/internal/utils/exception"
	"github.com/steppbol/activity-manager/internal/utils/mapper"
)

type ActivityRouter struct {
	activityService *services.ActivityService
}

func NewActivityRouter(r *gin.Engine, as *services.ActivityService) {
	dr := ActivityRouter{
		activityService: as,
	}

	api := r.Group("/api/v1/activity-manager")

	api.POST("/activities", dr.Create)
	api.PUT("/activities/:id", dr.Update)
	api.GET("/activities", dr.FindAllByUserID)
	api.GET("/activities/tags", dr.FindAllByTags)
	api.DELETE("/activities", dr.Delete)
}

func (ar ActivityRouter) Create(c *gin.Context) {
	var input dtos.ActivityDTO

	err := c.ShouldBindJSON(&input)
	if err != nil {
		dtos.CreateResponse(c, http.StatusBadRequest, exception.BadRequest, nil)
		return
	}

	activity := ar.activityService.Create(input.Title, input.Description, input.Content, input.DateID, input.TagIDs)
	if activity == nil {
		dtos.CreateResponse(c, http.StatusConflict, exception.Conflict, nil)
		return
	}

	dtos.CreateResponse(c, http.StatusOK, exception.Success, activity)
}

func (ar ActivityRouter) Update(c *gin.Context) {
	var input dtos.ActivityUpdateRequest

	err := c.ShouldBindJSON(&input)
	if err != nil {
		dtos.CreateResponse(c, http.StatusBadRequest, exception.BadRequest, nil)
		return
	}

	id := c.Param("id")

	cId, err := strconv.Atoi(id)
	if err != nil {
		dtos.CreateResponse(c, http.StatusBadRequest, exception.BadRequest, nil)
		return
	}

	user := ar.activityService.Update(uint(cId), *mapper.ActivityUpdateRequestToMap(input))

	dtos.CreateResponse(c, http.StatusOK, exception.Success, user)
}

func (ar ActivityRouter) FindAllByUserID(c *gin.Context) {
	var input dtos.FindByUserIDRequest

	err := c.ShouldBindJSON(&input)
	if err != nil {
		dtos.CreateResponse(c, http.StatusBadRequest, exception.BadRequest, nil)
		return
	}

	activity := ar.activityService.FindAllByUserID(input.UserID)

	dtos.CreateResponse(c, http.StatusOK, exception.Success, activity)
}

func (ar ActivityRouter) FindAllByTags(c *gin.Context) {
	var input dtos.FindByTagsRequest

	err := c.ShouldBindJSON(&input)
	if err != nil {
		dtos.CreateResponse(c, http.StatusBadRequest, exception.BadRequest, nil)
		return
	}

	activity := ar.activityService.FindAllByTags(input.UserID, input.TagIDs)

	dtos.CreateResponse(c, http.StatusOK, exception.Success, activity)
}

func (ar ActivityRouter) Delete(c *gin.Context) {
	id := c.Param("id")

	cId, err := strconv.Atoi(id)
	if err != nil {
		dtos.CreateResponse(c, http.StatusBadRequest, exception.BadRequest, nil)
		return
	}

	ar.activityService.DeleteByID(uint(cId))

	dtos.CreateResponse(c, http.StatusOK, exception.Success, nil)
}
