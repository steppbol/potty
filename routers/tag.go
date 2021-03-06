package routers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/steppbol/activity-manager/dtos"
	"github.com/steppbol/activity-manager/services"
	"github.com/steppbol/activity-manager/utils/exception"
)

type TagRouter struct {
	tagService *services.TagService
}

func NewTagRouter(r *gin.Engine, ts *services.TagService) {
	tr := TagRouter{
		tagService: ts,
	}

	api := r.Group("/api/v1/activity-manager")

	api.POST("/tags", tr.Create)
	api.PUT("/tags/:id", tr.Update)
	api.GET("/tags", tr.FindAll)
	api.GET("/tags/:id", tr.FindByID)
	api.DELETE("/tags/:id", tr.Delete)
}

func (tr TagRouter) Create(c *gin.Context) {
	var input dtos.TagDTO

	err := c.ShouldBindJSON(&input)
	if err != nil {
		dtos.CreateResponse(c, http.StatusBadRequest, exception.BadRequest, nil)
		return
	}

	tag := tr.tagService.Create(input.Name)
	if tag == nil {
		dtos.CreateResponse(c, http.StatusConflict, exception.Conflict, nil)
		return
	}

	dtos.CreateResponse(c, http.StatusOK, exception.Success, tag)
}

func (tr TagRouter) Update(c *gin.Context) {
	var input dtos.TagDTO

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

	tag := tr.tagService.Update(uint(cId), input.Name)
	if tag == nil {
		dtos.CreateResponse(c, http.StatusConflict, exception.Conflict, nil)
		return
	}

	dtos.CreateResponse(c, http.StatusOK, exception.Success, tag)
}

func (tr TagRouter) FindAll(c *gin.Context) {
	tags := tr.tagService.FindAll()
	dtos.CreateResponse(c, http.StatusOK, exception.Success, tags)
}

func (tr TagRouter) FindByID(c *gin.Context) {
	id := c.Param("id")

	cId, err := strconv.Atoi(id)
	if err != nil {
		dtos.CreateResponse(c, http.StatusBadRequest, exception.BadRequest, nil)
		return
	}

	tag, err := tr.tagService.FindByID(uint(cId))
	if err != nil {
		dtos.CreateResponse(c, http.StatusNotFound, exception.NotFound, nil)
		return
	}

	dtos.CreateResponse(c, http.StatusOK, exception.Success, tag)
}

func (tr TagRouter) Delete(c *gin.Context) {
	id := c.Param("id")

	cId, err := strconv.Atoi(id)
	if err != nil {
		dtos.CreateResponse(c, http.StatusBadRequest, exception.BadRequest, nil)
		return
	}

	tr.tagService.DeleteByID(uint(cId))

	dtos.CreateResponse(c, http.StatusOK, exception.Success, nil)
}
