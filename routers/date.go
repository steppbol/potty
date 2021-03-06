package routers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/steppbol/activity-manager/dtos"
	"github.com/steppbol/activity-manager/services"
	"github.com/steppbol/activity-manager/utils/exception"
	"github.com/steppbol/activity-manager/utils/mapper"
)

type DateRouter struct {
	dateService *services.DateService
}

func NewDateRouter(r *gin.Engine, ds *services.DateService) {
	dr := DateRouter{
		dateService: ds,
	}

	api := r.Group("/api/v1/activity-manager")

	api.POST("/dates", dr.Create)
	api.PUT("/dates/:id", dr.Update)
	api.GET("/dates/:id", dr.FindByID)
	api.GET("/users/:id/dates", dr.FindByUserID)
	api.GET("/users/:id/dates/:id2/activities", dr.FindAllActivities)
	api.DELETE("/dates", dr.Delete)
}

func (dr DateRouter) Create(c *gin.Context) {
	var input dtos.DateDTO

	err := c.ShouldBindJSON(&input)
	if err != nil {
		dtos.CreateResponse(c, http.StatusBadRequest, exception.BadRequest, nil)
		return
	}

	date := dr.dateService.Create(input.Time, input.UserID, input.Note)
	if date == nil {
		dtos.CreateResponse(c, http.StatusConflict, exception.Conflict, nil)
		return
	}

	dtos.CreateResponse(c, http.StatusOK, exception.Success, date)
}

func (dr DateRouter) Update(c *gin.Context) {
	var input dtos.DateUpdateRequest

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

	user := dr.dateService.Update(uint(cId), *mapper.DateUpdateRequestToMap(input))

	dtos.CreateResponse(c, http.StatusOK, exception.Success, user)
}

func (dr DateRouter) FindByID(c *gin.Context) {
	id := c.Param("id")

	cId, err := strconv.Atoi(id)
	if err != nil {
		dtos.CreateResponse(c, http.StatusBadRequest, exception.BadRequest, nil)
		return
	}

	date, err := dr.dateService.FindByID(uint(cId))

	if err != nil {
		dtos.CreateResponse(c, http.StatusNotFound, exception.NotFound, nil)
		return
	}

	dtos.CreateResponse(c, http.StatusOK, exception.Success, date)
}

func (dr DateRouter) FindAllActivities(c *gin.Context) {
	uId := c.Param("id")

	cuId, err := strconv.Atoi(uId)
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

	date := dr.dateService.FindAllActivities(uint(cId), uint(cuId))

	dtos.CreateResponse(c, http.StatusOK, exception.Success, date)
}

func (dr DateRouter) FindByUserID(c *gin.Context) {
	id := c.Param("id")

	cId, err := strconv.Atoi(id)
	if err != nil {
		dtos.CreateResponse(c, http.StatusBadRequest, exception.BadRequest, nil)
		return
	}

	date := dr.dateService.FindAllByUserID(uint(cId))

	dtos.CreateResponse(c, http.StatusOK, exception.Success, date)
}

func (dr DateRouter) Delete(c *gin.Context) {
	id := c.Param("id")

	cId, err := strconv.Atoi(id)
	if err != nil {
		dtos.CreateResponse(c, http.StatusBadRequest, exception.BadRequest, nil)
		return
	}

	dr.dateService.DeleteByID(uint(cId))

	dtos.CreateResponse(c, http.StatusOK, exception.Success, nil)
}
