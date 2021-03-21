package cache

import (
	"fmt"

	"github.com/go-redis/redis/v8"

	"github.com/steppbol/activity-manager/configs"
)

type RedisCache struct {
	Client *redis.Client
	config *configs.Cache
}

func NewRedisCache(conf *configs.Cache) *RedisCache {
	c := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		Password: conf.Password,
		DB:       0,
	})

	return &RedisCache{
		Client: c,
		config: conf,
	}
}
