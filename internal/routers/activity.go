package routers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/steppbol/activity-manager/internal/api"
	"github.com/steppbol/activity-manager/internal/dtos"
	"github.com/steppbol/activity-manager/internal/middleware"
	"github.com/steppbol/activity-manager/internal/utils/exception"
	"github.com/steppbol/activity-manager/internal/utils/mapper"
)

type ActivityRouter struct {
	baseAPI *api.BaseAPI
}

func NewActivityRouter(r *gin.Engine, ba *api.BaseAPI, jm *middleware.JWTMiddleware) {
	dr := ActivityRouter{
		baseAPI: ba,
	}

	routers := r.Group("/api/v1/activity-manager")

	routers.Use(jm.JWT())
	{
		routers.POST("/activities", dr.Create)
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

	activity := ar.baseAPI.ActivityService.Create(input.Title, input.Description, input.Content, input.DateID, input.TagIDs)
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

	user := ar.baseAPI.ActivityService.Update(uint(cId), *mapper.ActivityUpdateRequestToMap(input))

	dtos.CreateJSONResponse(c, http.StatusOK, exception.Success, user)
}

func (ar ActivityRouter) FindAllByUserID(c *gin.Context) {
	var input dtos.UserIDRequest

	err := c.ShouldBindJSON(&input)
	if err != nil {
		dtos.CreateJSONResponse(c, http.StatusBadRequest, exception.BadRequest, nil)
		return
	}

	activity := ar.baseAPI.ActivityService.FindAllByUserID(input.UserID)

	dtos.CreateJSONResponse(c, http.StatusOK, exception.Success, activity)
}

func (ar ActivityRouter) FindAllByTags(c *gin.Context) {
	var input dtos.FindByTagsRequest

	err := c.ShouldBindJSON(&input)
	if err != nil {
		dtos.CreateJSONResponse(c, http.StatusBadRequest, exception.BadRequest, nil)
		return
	}

	activity := ar.baseAPI.ActivityService.FindAllByTags(input.UserID, input.TagIDs)

	dtos.CreateJSONResponse(c, http.StatusOK, exception.Success, activity)
}

func (ar ActivityRouter) Delete(c *gin.Context) {
	id := c.Param("id")

	cId, err := strconv.Atoi(id)
	if err != nil {
		dtos.CreateJSONResponse(c, http.StatusBadRequest, exception.BadRequest, nil)
		return
	}

	ar.baseAPI.ActivityService.DeleteByID(uint(cId))

	dtos.CreateJSONResponse(c, http.StatusOK, exception.Success, nil)
}
