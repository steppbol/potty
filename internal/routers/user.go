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

type UserRouter struct {
	baseAPI *api.BaseAPI
}

func NewUserRouter(r *gin.Engine, ba *api.BaseAPI, jm *middleware.JWTMiddleware) {
	ur := UserRouter{
		baseAPI: ba,
	}

	routers := r.Group("/api/v1/activity-manager")

	routers.POST("/users", ur.Create)

	routers.Use(jm.JWT())
	{
		routers.PUT("/users/:id", ur.Update)
		routers.GET("/users/:id", ur.FindByID)
		routers.DELETE("/users/:id", ur.Delete)
	}
}

func (ur UserRouter) Create(c *gin.Context) {
	var input dtos.UserDTO

	err := c.ShouldBindJSON(&input)
	if err != nil {
		dtos.CreateJSONResponse(c, http.StatusBadRequest, exception.BadRequest, nil)
		return
	}

	user := ur.baseAPI.UserService.Create(input.Username, input.Password, input.Email)
	if user == nil {
		dtos.CreateJSONResponse(c, http.StatusConflict, exception.Conflict, nil)
		return
	}

	dtos.CreateJSONResponse(c, http.StatusOK, exception.Success, user)
}

func (ur UserRouter) Update(c *gin.Context) {
	var input dtos.UserUpdateRequest

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

	user := ur.baseAPI.UserService.Update(uint(cId), *mapper.UserUpdateRequestToMap(input))

	dtos.CreateJSONResponse(c, http.StatusOK, exception.Success, user)
}

func (ur UserRouter) FindByID(c *gin.Context) {
	id := c.Param("id")

	cId, err := strconv.Atoi(id)
	if err != nil {
		dtos.CreateJSONResponse(c, http.StatusBadRequest, exception.BadRequest, nil)
		return
	}

	user, err := ur.baseAPI.UserService.FindByID(uint(cId))

	if err != nil {
		dtos.CreateJSONResponse(c, http.StatusNotFound, exception.NotFound, nil)
		return
	}

	dtos.CreateJSONResponse(c, http.StatusOK, exception.Success, user)
}

func (ur UserRouter) Delete(c *gin.Context) {
	id := c.Param("id")

	cId, err := strconv.Atoi(id)
	if err != nil {
		dtos.CreateJSONResponse(c, http.StatusBadRequest, exception.BadRequest, nil)
		return
	}

	ur.baseAPI.UserService.DeleteByID(uint(cId))

	dtos.CreateJSONResponse(c, http.StatusOK, exception.Success, nil)
}
