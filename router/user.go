package router

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/steppbol/activity-manager/dto"
	"github.com/steppbol/activity-manager/service"
	"github.com/steppbol/activity-manager/util/exception"
	"github.com/steppbol/activity-manager/util/mapper"
)

type UserRouter struct {
	userService *service.UserService
}

func NewUserRouter(r *gin.Engine, us *service.UserService) {
	ur := UserRouter{
		userService: us,
	}

	api := r.Group("/api/v1/activity-manager")

	api.POST("/users", ur.Create)
	api.PUT("/users/:id", ur.Update)
	api.GET("/users/:id", ur.FindByID)
	api.DELETE("/users/:id", ur.Delete)
}

func (ur UserRouter) Create(c *gin.Context) {
	var input dto.UserDTO

	err := c.ShouldBindJSON(&input)
	if err != nil {
		dto.CreateResponse(c, http.StatusBadRequest, exception.BadRequest, nil)
		return
	}

	tag := ur.userService.Create(input.Username, input.Password)
	if tag == nil {
		dto.CreateResponse(c, http.StatusConflict, exception.Conflict, nil)
		return
	}

	dto.CreateResponse(c, http.StatusOK, exception.Success, tag)
}

func (ur UserRouter) Update(c *gin.Context) {
	var input dto.UserUpdateRequest

	err := c.ShouldBindJSON(&input)
	if err != nil {
		dto.CreateResponse(c, http.StatusBadRequest, exception.BadRequest, nil)
		return
	}

	id := c.Param("id")

	cId, err := strconv.Atoi(id)
	if err != nil {
		dto.CreateResponse(c, http.StatusBadRequest, exception.BadRequest, nil)
		return
	}

	user := ur.userService.Update(uint(cId), *mapper.UserUpdateRequestToMap(input))

	dto.CreateResponse(c, http.StatusOK, exception.Success, user)
}

func (ur UserRouter) FindByID(c *gin.Context) {
	id := c.Param("id")

	cId, err := strconv.Atoi(id)
	if err != nil {
		dto.CreateResponse(c, http.StatusBadRequest, exception.BadRequest, nil)
		return
	}

	user, err := ur.userService.FindByID(uint(cId))

	if err != nil {
		dto.CreateResponse(c, http.StatusNotFound, exception.NotFound, nil)
		return
	}

	dto.CreateResponse(c, http.StatusOK, exception.Success, user)
}

func (ur UserRouter) Delete(c *gin.Context) {
	id := c.Param("id")

	cId, err := strconv.Atoi(id)
	if err != nil {
		dto.CreateResponse(c, http.StatusBadRequest, exception.BadRequest, nil)
		return
	}

	ur.userService.DeleteByID(uint(cId))

	dto.CreateResponse(c, http.StatusOK, exception.Success, nil)
}
