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

type DateRouter struct {
	dateService *services.DateService
}

func NewDateRouter(r *gin.Engine, ds *services.DateService) {
	dr := DateRouter{
		dateService: ds,
	}

	api := r.Group("/api/v1/activity-manager")

	api.POST("/dates", dr.Create)
	api.POST("/dates/export", dr.ExportToXLSX)
	api.PUT("/dates/:id", dr.Update)
	api.GET("/dates/:id", dr.FindByID)
	api.GET("/dates", dr.FindAllByUserID)
	api.DELETE("/dates", dr.Delete)
}

func (dr DateRouter) Create(c *gin.Context) {
	var input dtos.DateDTO

	err := c.ShouldBindJSON(&input)
	if err != nil {
		dtos.CreateJSONResponse(c, http.StatusBadRequest, exception.BadRequest, nil)
		return
	}

	date := dr.dateService.Create(input.Time, input.UserID, input.Note)
	if date == nil {
		dtos.CreateJSONResponse(c, http.StatusConflict, exception.Conflict, nil)
		return
	}

	dtos.CreateJSONResponse(c, http.StatusOK, exception.Success, date)
}

func (dr DateRouter) ExportToXLSX(c *gin.Context) {
	var input dtos.UserIDRequest

	err := c.ShouldBindJSON(&input)
	if err != nil {
		dtos.CreateJSONResponse(c, http.StatusBadRequest, exception.BadRequest, nil)
		return
	}

	file, err := dr.dateService.ExportToXLSX(input.UserID)
	if err != nil {
		dtos.CreateJSONResponse(c, http.StatusInternalServerError, exception.InternalServerError, nil)
		return
	}

	dtos.CreateBinResponse(c, http.StatusInternalServerError, file)
}

func (dr DateRouter) Update(c *gin.Context) {
	var input dtos.DateUpdateRequest

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

	user := dr.dateService.Update(uint(cId), *mapper.DateUpdateRequestToMap(input))

	dtos.CreateJSONResponse(c, http.StatusOK, exception.Success, user)
}

func (dr DateRouter) FindByID(c *gin.Context) {
	id := c.Param("id")

	cId, err := strconv.Atoi(id)
	if err != nil {
		dtos.CreateJSONResponse(c, http.StatusBadRequest, exception.BadRequest, nil)
		return
	}

	date, err := dr.dateService.FindByID(uint(cId))

	if err != nil {
		dtos.CreateJSONResponse(c, http.StatusNotFound, exception.NotFound, nil)
		return
	}

	dtos.CreateJSONResponse(c, http.StatusOK, exception.Success, date)
}

func (dr DateRouter) FindAllByUserID(c *gin.Context) {
	var input dtos.UserIDRequest

	err := c.ShouldBindJSON(&input)
	if err != nil {
		dtos.CreateJSONResponse(c, http.StatusBadRequest, exception.BadRequest, nil)
		return
	}

	date := dr.dateService.FindAllByUserID(input.UserID)

	dtos.CreateJSONResponse(c, http.StatusOK, exception.Success, date)
}

func (dr DateRouter) Delete(c *gin.Context) {
	id := c.Param("id")

	cId, err := strconv.Atoi(id)
	if err != nil {
		dtos.CreateJSONResponse(c, http.StatusBadRequest, exception.BadRequest, nil)
		return
	}

	dr.dateService.DeleteByID(uint(cId))

	dtos.CreateJSONResponse(c, http.StatusOK, exception.Success, nil)
}
