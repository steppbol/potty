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

type DateRouter struct {
	baseAPI *api.BaseAPI
}

func NewDateRouter(r *gin.Engine, ba *api.BaseAPI, jm *middleware.JWTMiddleware) {
	dr := DateRouter{
		baseAPI: ba,
	}

	routers := r.Group("/api/v1/activity-manager")

	routers.Use(jm.JWT())
	{
		routers.POST("/dates", dr.Create)
		routers.POST("/dates/export", dr.ExportToXLSX)
		routers.POST("/dates/import/:id", dr.ImportFromXLSX)
		routers.PUT("/dates/:id", dr.Update)
		routers.GET("/dates/:id", dr.FindByID)
		routers.GET("/dates", dr.FindAllByUserID)
		routers.DELETE("/dates", dr.Delete)
	}
}

func (dr DateRouter) Create(c *gin.Context) {
	var input dtos.DateDTO

	err := c.ShouldBindJSON(&input)
	if err != nil {
		dtos.CreateJSONResponse(c, http.StatusBadRequest, exception.BadRequest, nil)
		return
	}

	date := dr.baseAPI.DateService.Create(input.Time, input.UserID, input.Note)
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

	path, err := dr.baseAPI.ExportToXLSX(input.UserID)
	if err != nil {
		dtos.CreateJSONResponse(c, http.StatusInternalServerError, exception.InternalServerError, nil)
		return
	}

	dtos.CreateBinResponse(c, path)

	_ = dr.baseAPI.DeleteStaticData(path)
}

func (dr DateRouter) ImportFromXLSX(c *gin.Context) {
	id := c.Param("id")

	cId, err := strconv.Atoi(id)
	if err != nil {
		dtos.CreateJSONResponse(c, http.StatusBadRequest, exception.BadRequest, nil)
		return
	}

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		dtos.CreateJSONResponse(c, http.StatusInternalServerError, exception.InternalServerError, nil)
		return
	}

	err = dr.baseAPI.ImportFromXLSX(uint(cId), file)
	if err != nil {
		dtos.CreateJSONResponse(c, http.StatusInternalServerError, exception.InternalServerError, nil)
		return
	}

	dtos.CreateJSONResponse(c, http.StatusOK, exception.Success, nil)
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

	user := dr.baseAPI.DateService.Update(uint(cId), *mapper.DateUpdateRequestToMap(input))

	dtos.CreateJSONResponse(c, http.StatusOK, exception.Success, user)
}

func (dr DateRouter) FindByID(c *gin.Context) {
	id := c.Param("id")

	cId, err := strconv.Atoi(id)
	if err != nil {
		dtos.CreateJSONResponse(c, http.StatusBadRequest, exception.BadRequest, nil)
		return
	}

	date, err := dr.baseAPI.DateService.FindByID(uint(cId))

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

	date := dr.baseAPI.DateService.FindAllByUserID(input.UserID)

	dtos.CreateJSONResponse(c, http.StatusOK, exception.Success, date)
}

func (dr DateRouter) Delete(c *gin.Context) {
	id := c.Param("id")

	cId, err := strconv.Atoi(id)
	if err != nil {
		dtos.CreateJSONResponse(c, http.StatusBadRequest, exception.BadRequest, nil)
		return
	}

	dr.baseAPI.DateService.DeleteByID(uint(cId))

	dtos.CreateJSONResponse(c, http.StatusOK, exception.Success, nil)
}
