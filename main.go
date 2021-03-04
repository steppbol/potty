package main

import (
	"github.com/gin-gonic/gin"

	"github.com/steppbol/activity-manager/repository"
	"github.com/steppbol/activity-manager/router"
	"github.com/steppbol/activity-manager/service"
)

func main() {
	br, err := repository.Setup()
	if err != nil {
		panic(err)
	}

	_, err = repository.NewActivityRepository(br)
	if err != nil {
		panic(err)
	}

	dr, err := repository.NewDateRepository(br)
	if err != nil {
		panic(err)
	}

	tr, err := repository.NewTagRepository(br)
	if err != nil {
		panic(err)
	}

	ur, err := repository.NewUserRepository(br)
	if err != nil {
		panic(err)
	}

	ts, err := service.NewTagService(tr)
	if err != nil {
		panic(err)
	}

	us, err := service.NewUserService(ur)
	if err != nil {
		panic(err)
	}

	ds, err := service.NewDateService(dr)
	if err != nil {
		panic(err)
	}

	r := gin.Default()
	router.NewTagRouter(r, ts)
	router.NewUserRouter(r, us)
	router.NewDateRouter(r, ds)

	err = r.Run()
	if err != nil {
		panic(err)
	}
}
