package repositories

import (
	"fmt"
	"net/url"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/steppbol/activity-manager/configs"
	"github.com/steppbol/activity-manager/internal/models"
)

type BaseRepository struct {
	database *gorm.DB
	config   *configs.Database
}

func Setup(conf *configs.Database) (*BaseRepository, error) {
	dsn := url.URL{
		User:     url.UserPassword(conf.User, conf.Password),
		Scheme:   conf.Schema,
		Host:     fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		Path:     conf.Name,
		RawQuery: (&url.Values{"sslmode": []string{conf.SSLMode}}).Encode(),
	}

	time.Sleep(5 * time.Second)

	db, err := gorm.Open(postgres.Open(dsn.String()), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if conf.GenerateSchema {
		err = db.AutoMigrate(&models.Tag{})
		if err != nil {
			return nil, err
		}

		err = db.AutoMigrate(&models.Activity{})
		if err != nil {
			return nil, err
		}

		err = db.AutoMigrate(&models.Date{})
		if err != nil {
			return nil, err
		}

		err = db.AutoMigrate(&models.User{})
		if err != nil {
			return nil, err
		}
	}

	return &BaseRepository{
		database: db,
		config:   conf,
	}, nil
}
