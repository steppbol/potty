package routers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/steppbol/activity-manager/internal/api"
	"github.com/steppbol/activity-manager/internal/dtos"
	"github.com/steppbol/activity-manager/internal/utils/exception"
)

type TagRouter struct {
	baseAPI *api.BaseAPI
}

func NewTagRouter(r *gin.Engine, ba *api.BaseAPI) {
	tr := TagRouter{
		baseAPI: ba,
	}

	routers := r.Group("/api/v1/activity-manager")

	routers.POST("/tags", tr.Create)
	routers.PUT("/tags/:id", tr.Update)
	routers.GET("/tags", tr.FindAll)
	routers.GET("/tags/:id", tr.FindByID)
	routers.DELETE("/tags/:id", tr.Delete)
}

func (tr TagRouter) Create(c *gin.Context) {
	var input dtos.TagDTO

	err := c.ShouldBindJSON(&input)
	if err != nil {
		dtos.CreateJSONResponse(c, http.StatusBadRequest, exception.BadRequest, nil)
		return
	}

	tag := tr.baseAPI.TagService.Create(input.Name)
	if tag == nil {
		dtos.CreateJSONResponse(c, http.StatusConflict, exception.Conflict, nil)
		return
	}

	dtos.CreateJSONResponse(c, http.StatusOK, exception.Success, tag)
}

func (tr TagRouter) Update(c *gin.Context) {
	var input dtos.TagDTO

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

	tag := tr.baseAPI.TagService.Update(uint(cId), input.Name)
	if tag == nil {
		dtos.CreateJSONResponse(c, http.StatusConflict, exception.Conflict, nil)
		return
	}

	dtos.CreateJSONResponse(c, http.StatusOK, exception.Success, tag)
}

func (tr TagRouter) FindAll(c *gin.Context) {
	tags := tr.baseAPI.TagService.FindAll()
	dtos.CreateJSONResponse(c, http.StatusOK, exception.Success, tags)
}

func (tr TagRouter) FindByID(c *gin.Context) {
	id := c.Param("id")

	cId, err := strconv.Atoi(id)
	if err != nil {
		dtos.CreateJSONResponse(c, http.StatusBadRequest, exception.BadRequest, nil)
		return
	}

	tag, err := tr.baseAPI.TagService.FindByID(uint(cId))
	if err != nil {
		dtos.CreateJSONResponse(c, http.StatusNotFound, exception.NotFound, nil)
		return
	}

	dtos.CreateJSONResponse(c, http.StatusOK, exception.Success, tag)
}

func (tr TagRouter) Delete(c *gin.Context) {
	id := c.Param("id")

	cId, err := strconv.Atoi(id)
	if err != nil {
		dtos.CreateJSONResponse(c, http.StatusBadRequest, exception.BadRequest, nil)
		return
	}

	tr.baseAPI.TagService.DeleteByID(uint(cId))

	dtos.CreateJSONResponse(c, http.StatusOK, exception.Success, nil)
}
