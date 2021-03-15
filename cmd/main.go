package main

import (
	"github.com/gin-gonic/gin"

	"github.com/steppbol/activity-manager/configs"
	"github.com/steppbol/activity-manager/internal/api"
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
	as := services.NewActivityService(ts, ds, ar)

	ba := api.NewBaseAPI(ts, as, ds, us)

	r := gin.Default()

	routers.NewUserRouter(r, ba)
	routers.NewTagRouter(r, ba)
	routers.NewDateRouter(r, ba)
	routers.NewActivityRouter(r, ba)

	err = r.Run()
	if err != nil {
		panic(err)
	}
}
