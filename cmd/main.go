package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/steppbol/activity-manager/configs"
	"github.com/steppbol/activity-manager/internal/api"
	"github.com/steppbol/activity-manager/internal/cache"
	"github.com/steppbol/activity-manager/internal/middleware"
	"github.com/steppbol/activity-manager/internal/repositories"
	"github.com/steppbol/activity-manager/internal/routers"
	"github.com/steppbol/activity-manager/internal/services"
)

func main() {
	c, err := configs.Setup()
	if err != nil {
		panic(err)
	}

	br, err := repositories.Setup(&c.Database)
	if err != nil {
		panic(err)
	}

	ur := repositories.NewUserRepository(br)
	tr := repositories.NewTagRepository(br)
	ar := repositories.NewActivityRepository(br)
	dr := repositories.NewDateRepository(br)

	xs := services.NewXLSXService(&c.Application)
	us := services.NewUserService(ur)
	ts := services.NewTagService(tr)
	ds := services.NewDateService(us, xs, dr)
	as := services.NewActivityService(ts, ds, us, ar)

	rc := cache.NewRedisCache(&c.Cache)

	auths := services.NewAuthenticationService(rc)

	ba := api.NewXLSXBaseAPI(ts, as, ds)

	jm, err := middleware.NewJWTMiddleware(&c.Security)
	if err != nil {
		panic(err)
	}

	gin.SetMode(c.Server.Mode)

	r := gin.New()

	routers.NewUserRouter(r, us, jm)
	routers.NewTagRouter(r, ts, jm)
	routers.NewDateRouter(r, ba, ds, jm)
	routers.NewActivityRouter(r, as, jm)
	routers.NewAuthenticationRouter(r, us, auths, jm)

	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", c.Server.Port),
		Handler:        r,
		MaxHeaderBytes: 1 << 20,
	}

	log.Info("application is started. Port : ", c.Server.Port)

	err = server.ListenAndServe()

	if err != nil {
		panic(err)
	}
}
