package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/steppbol/activity-manager/dto"
	"github.com/steppbol/activity-manager/service"
	"github.com/steppbol/activity-manager/util/exception"
)

type DateRouter struct {
	dateService *service.DateService
}

func NewDateRouter(r *gin.Engine, ds *service.DateService) {
	dr := DateRouter{
		dateService: ds,
	}

	api := r.Group("/api/v1/activity-manager")

	api.POST("/dates", dr.Create)
}

func (dr DateRouter) Create(c *gin.Context) {
	var input dto.DateDTO

	err := c.ShouldBindJSON(&input)
	if err != nil {
		dto.CreateResponse(c, http.StatusBadRequest, exception.BadRequest, nil)
		return
	}

	date := dr.dateService.Create(input.Time, input.UserID, input.Note)
	if date == nil {
		dto.CreateResponse(c, http.StatusConflict, exception.Conflict, nil)
		return
	}

	dto.CreateResponse(c, http.StatusOK, exception.Success, date)
}
