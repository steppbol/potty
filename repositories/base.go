package repositories

import (
	"fmt"
	"net/url"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/steppbol/activity-manager/configs"
	"github.com/steppbol/activity-manager/models"
)

type BaseRepository struct {
	database *gorm.DB
}

func Setup() (*BaseRepository, error) {
	conf, err := configs.GetDatabaseConfig()

	if err != nil {
		return nil, err
	}

	dsn := url.URL{
		User:     url.UserPassword(conf.User, conf.Password),
		Scheme:   conf.Scheme,
		Host:     fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		Path:     conf.Name,
		RawQuery: (&url.Values{"sslmode": []string{"disable"}}).Encode(),
	}

	db, err := gorm.Open(postgres.Open(dsn.String()), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&models.Tag{})
	db.AutoMigrate(&models.Activity{})
	db.AutoMigrate(&models.Date{})
	db.AutoMigrate(&models.User{})

	return &BaseRepository{
		database: db,
	}, nil
}
