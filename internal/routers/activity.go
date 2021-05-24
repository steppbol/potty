package routers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/steppbol/activity-manager/internal/dtos"
	"github.com/steppbol/activity-manager/internal/middleware"
	"github.com/steppbol/activity-manager/internal/services"
	"github.com/steppbol/activity-manager/internal/utils/exception"
	"github.com/steppbol/activity-manager/internal/utils/mapper"
)

type ActivityRouter struct {
	activityService *services.ActivityService
}

func NewActivityRouter(r *gin.Engine, as *services.ActivityService, jm *middleware.JWTMiddleware) {
	dr := ActivityRouter{
		activityService: as,
	}

	routers := r.Group("/api/v1/activity-manager")

	routers.Use(jm.CORS())

	routers.Use(jm.JWT())
	{
		routers.POST("/activities", dr.Create)
		routers.POST("/activities/strict", dr.CreateWithDateID)
		routers.PUT("/activities/:id", dr.Update)
		routers.GET("/activities", dr.FindAllByUserID)
		routers.GET("/activities/tags", dr.FindAllByTags)
		routers.DELETE("/activities", dr.Delete)
	}
}

func (ar ActivityRouter) Create(c *gin.Context) {
	var input dtos.ActivityDTO

	err := c.ShouldBindJSON(&input)
	if err != nil {
		dtos.CreateJSONResponse(c, http.StatusBadRequest, exception.BadRequest, nil)
		return
	}

	cActivity := mapper.ActivityDTOToActivity(input)

	activity := ar.activityService.Create(*cActivity, input.Username, input.Date, input.TagIDs)
	if activity == nil {
		dtos.CreateJSONResponse(c, http.StatusConflict, exception.Conflict, nil)
		return
	}

	dtos.CreateJSONResponse(c, http.StatusOK, exception.Success, activity)
}

func (ar ActivityRouter) CreateWithDateID(c *gin.Context) {
	var input dtos.ActivityWithDateIDDTO

	err := c.ShouldBindJSON(&input)
	if err != nil {
		dtos.CreateJSONResponse(c, http.StatusBadRequest, exception.BadRequest, nil)
		return
	}

	cActivity := mapper.ActivityWithDateIDDTOToActivity(input)

	activity := ar.activityService.CreateWithDateID(*cActivity, input.TagIDs)
	if activity == nil {
		dtos.CreateJSONResponse(c, http.StatusConflict, exception.Conflict, nil)
		return
	}

	dtos.CreateJSONResponse(c, http.StatusOK, exception.Success, activity)
}

func (ar ActivityRouter) Update(c *gin.Context) {
	var input dtos.ActivityUpdateRequest

	err := c.ShouldBindJSON(&input)
	if err != nil {
		dtos.CreateJSONResponse(c, http.StatusBadRequest, exception.BadRequest, nil)
		return
	}

	id := c.Param("id")

	cId, err := strconv.Atoi(id)
	if err != nil {
		dtos.CreateJSONResponse(c, http.StatusBadRequest, exception.BadRequest, nil)
		return
	}

	user := ar.activityService.Update(uint(cId), *mapper.ActivityUpdateRequestToMap(input))

	dtos.CreateJSONResponse(c, http.StatusOK, exception.Success, user)
}

func (ar ActivityRouter) FindAllByUserID(c *gin.Context) {
	var input dtos.UserIDRequest

	err := c.ShouldBindJSON(&input)
	if err != nil {
		dtos.CreateJSONResponse(c, http.StatusBadRequest, exception.BadRequest, nil)
		return
	}

	activity := ar.activityService.FindAllByUserID(input.UserID)

	dtos.CreateJSONResponse(c, http.StatusOK, exception.Success, activity)
}

func (ar ActivityRouter) FindAllByTags(c *gin.Context) {
	var input dtos.FindByTagsRequest

	err := c.ShouldBindJSON(&input)
	if err != nil {
		dtos.CreateJSONResponse(c, http.StatusBadRequest, exception.BadRequest, nil)
		return
	}

	activity := ar.activityService.FindAllByTags(input.UserID, input.TagIDs)

	dtos.CreateJSONResponse(c, http.StatusOK, exception.Success, activity)
}

func (ar ActivityRouter) Delete(c *gin.Context) {
	id := c.Param("id")

	cId, err := strconv.Atoi(id)
	if err != nil {
		dtos.CreateJSONResponse(c, http.StatusBadRequest, exception.BadRequest, nil)
		return
	}

	ar.activityService.DeleteByID(uint(cId))

	dtos.CreateJSONResponse(c, http.StatusOK, exception.Success, nil)
}
