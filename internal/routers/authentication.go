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
	userService           *services.UserService
	authenticationService *services.AuthenticationService
	jwtMiddleware         *middleware.JWTMiddleware
}

func NewAuthenticationRouter(r *gin.Engine, us *services.UserService, as *services.AuthenticationService, jm *middleware.JWTMiddleware) {
	dr := AuthenticationRouter{
		userService:           us,
		authenticationService: as,
		jwtMiddleware:         jm,
	}

	routers := r.Group("/api/v1/activity-manager/authentication")

	routers.Use(jm.CORS())

	routers.POST("/login", dr.Login)
	routers.POST("/refresh", dr.Refresh)

	routers.Use(jm.JWT())
	{
		routers.POST("/logout", dr.Logout)
	}
}

func (ar AuthenticationRouter) Login(c *gin.Context) {
	var input dtos.LoginRequest

	err := c.ShouldBindJSON(&input)
	if err != nil {
		dtos.CreateJSONResponse(c, http.StatusBadRequest, exception.BadRequest, nil)
		return
	}

	ok, u := ar.userService.CheckUser(input.Username, input.Password)
	if !ok {
		dtos.CreateJSONResponse(c, http.StatusUnauthorized, exception.Unauthorized, nil)
		return
	}

	token, err := ar.jwtMiddleware.GenerateToken(input.Username, input.Password, u.ID)

	err = ar.authenticationService.CreateAuthentication(u.ID, token)
	if err != nil {
		dtos.CreateJSONResponse(c, http.StatusInternalServerError, exception.InternalServerError, nil)
		return
	}

	tokens := map[string]string{
		"access_token":  token.AccessToken,
		"refresh_token": token.RefreshToken,
	}

	dtos.CreateJSONResponse(c, http.StatusOK, exception.Success, tokens)
}

func (ar AuthenticationRouter) Refresh(c *gin.Context) {
	var input dtos.RefreshRequest

	err := c.ShouldBindJSON(&input)
	if err != nil {
		dtos.CreateJSONResponse(c, http.StatusBadRequest, exception.BadRequest, nil)
		return
	}

	rd, err := ar.jwtMiddleware.ExtractRefreshTokenMetadata(input.RefreshToken)
	if err != nil {
		dtos.CreateJSONResponse(c, http.StatusUnauthorized, exception.Unauthorized, err)
		return
	}

	deleted, err := ar.authenticationService.DeleteAuthentication(rd.RefreshID)
	if err != nil || deleted == 0 {
		dtos.CreateJSONResponse(c, http.StatusUnauthorized, exception.Unauthorized, nil)
		return
	}

	u, err := ar.userService.FindByID(rd.UserID)
	if err != nil || deleted == 0 {
		dtos.CreateJSONResponse(c, http.StatusForbidden, exception.Forbidden, nil)
		return
	}

	ts, err := ar.jwtMiddleware.GenerateToken(u.Username, u.Password, u.ID)
	if err != nil {
		dtos.CreateJSONResponse(c, http.StatusForbidden, exception.Forbidden, nil)
		return
	}

	err = ar.authenticationService.CreateAuthentication(u.ID, ts)
	if err != nil {
		dtos.CreateJSONResponse(c, http.StatusForbidden, exception.Forbidden, nil)
		return
	}

	tokens := map[string]string{
		"access_token":  ts.AccessToken,
		"refresh_token": ts.RefreshToken,
	}

	dtos.CreateJSONResponse(c, http.StatusOK, exception.Success, tokens)
}

func (ar AuthenticationRouter) Logout(c *gin.Context) {
	metadata, err := ar.jwtMiddleware.ExtractAccessTokenMetadata(c)

	if err != nil {
		dtos.CreateJSONResponse(c, http.StatusUnauthorized, exception.Unauthorized, nil)
		return
	}

	ok, err := ar.authenticationService.DeleteTokens(metadata)
	if !ok || err != nil {
		dtos.CreateJSONResponse(c, http.StatusUnauthorized, exception.Unauthorized, nil)
		return
	}

	dtos.CreateJSONResponse(c, http.StatusOK, exception.Success, nil)
}
