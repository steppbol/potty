package main

import (
	"github.com/gin-gonic/gin"

	"github.com/steppbol/activity-manager/internal/repositories"
	"github.com/steppbol/activity-manager/internal/routers"
	"github.com/steppbol/activity-manager/internal/services"
)

func main() {
	br, err := repositories.Setup()
	if err != nil {
		panic(err)
	}

	ar, err := repositories.NewActivityRepository(br)
	if err != nil {
		panic(err)
	}

	dr, err := repositories.NewDateRepository(br)
	if err != nil {
		panic(err)
	}

	tr, err := repositories.NewTagRepository(br)
	if err != nil {
		panic(err)
	}

	ur, err := repositories.NewUserRepository(br)
	if err != nil {
		panic(err)
	}

	ts, err := services.NewTagService(tr)
	if err != nil {
		panic(err)
	}

	us, err := services.NewUserService(ur)
	if err != nil {
		panic(err)
	}

	ds, err := services.NewDateService(us, dr)
	if err != nil {
		panic(err)
	}

	as, err := services.NewActivityService(ds, ts, ar)
	if err != nil {
		panic(err)
	}

	r := gin.Default()
	routers.NewTagRouter(r, ts)
	routers.NewUserRouter(r, us)
	routers.NewDateRouter(r, ds)
	routers.NewActivityRouter(r, as)

	err = r.Run()
	if err != nil {
		panic(err)
	}
}
