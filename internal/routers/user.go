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

type UserRouter struct {
	userService *services.UserService
}

func NewUserRouter(r *gin.Engine, us *services.UserService, jm *middleware.JWTMiddleware) {
	ur := UserRouter{
		userService: us,
	}

	routers := r.Group("/api/v1/activity-manager")

	routers.POST("/users", ur.Create)

	routers.Use(jm.CORS())

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

	user := ur.userService.Create(input.Username, input.Password, input.Email)
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

	user := ur.userService.Update(uint(cId), *mapper.UserUpdateRequestToMap(input))

	dtos.CreateJSONResponse(c, http.StatusOK, exception.Success, user)
}

func (ur UserRouter) FindByID(c *gin.Context) {
	id := c.Param("id")

	cId, err := strconv.Atoi(id)
	if err != nil {
		dtos.CreateJSONResponse(c, http.StatusBadRequest, exception.BadRequest, nil)
		return
	}

	user, err := ur.userService.FindByID(uint(cId))

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

	ur.userService.DeleteByID(uint(cId))

	dtos.CreateJSONResponse(c, http.StatusOK, exception.Success, nil)
}
