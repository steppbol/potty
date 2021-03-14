package main

import (
	"github.com/gin-gonic/gin"

	"github.com/steppbol/activity-manager/configs"
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
	as := services.NewActivityService(ts, ar)
	ds := services.NewDateService(us, xs, dr)

	r := gin.Default()

	routers.NewUserRouter(r, us)
	routers.NewTagRouter(r, ts)
	routers.NewDateRouter(r, ds)
	routers.NewActivityRouter(r, as)

	err = r.Run()
	if err != nil {
		panic(err)
	}
}
