package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/steppbol/activity-manager/internal/dtos"
	"github.com/steppbol/activity-manager/internal/middleware"
	"github.com/steppbol/activity-manager/internal/services"
	"github.com/steppbol/activity-manager/internal/utils/exception"
)

type AuthenticationRouter struct {
	userService   *services.UserService
	jwtMiddleware *middleware.JWTMiddleware
}

func NewAuthenticationRouter(r *gin.Engine, us *services.UserService, jm *middleware.JWTMiddleware) {
	dr := AuthenticationRouter{
		userService:   us,
		jwtMiddleware: jm,
	}

	routers := r.Group("/api/v1/activity-manager/authentication")

	routers.POST("/login", dr.Login)

	routers.Use(jm.JWT())
	{
		routers.POST("/refresh", dr.Refresh)
	}
}

func (ar AuthenticationRouter) Login(c *gin.Context) {
	var input dtos.LoginRequest

	err := c.ShouldBindJSON(&input)
	if err != nil {
		dtos.CreateJSONResponse(c, http.StatusBadRequest, exception.BadRequest, nil)
		return
	}

	a := ar.userService.CheckUser(input.Username, input.Password)
	if !a {
		dtos.CreateJSONResponse(c, http.StatusUnauthorized, exception.Unauthorized, nil)
		return
	}

	token, err := ar.jwtMiddleware.GenerateToken(input.Username, input.Password)

	dtos.CreateJSONResponse(c, http.StatusOK, exception.Success, token)
}

func (ar AuthenticationRouter) Refresh(c *gin.Context) {
	var input dtos.UserIDRequest

	err := c.ShouldBindJSON(&input)
	if err != nil {
		dtos.CreateJSONResponse(c, http.StatusBadRequest, exception.BadRequest, nil)
		return
	}

	user, err := ar.userService.FindByID(input.UserID)
	if err != nil {
		dtos.CreateJSONResponse(c, http.StatusBadRequest, exception.BadRequest, nil)
		return
	}

	token, err := ar.jwtMiddleware.GenerateToken(user.Username, user.Password)
	if err != nil {
		dtos.CreateJSONResponse(c, http.StatusUnauthorized, exception.Unauthorized, nil)
		return
	}

	dtos.CreateJSONResponse(c, http.StatusOK, exception.Success, token)
}
