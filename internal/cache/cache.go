package cache

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"

	"github.com/steppbol/activity-manager/configs"
)

type RedisCache struct {
	connection *redis.Pool
	config     *configs.Cache
}

func NewRedisCache(conf *configs.Cache) *RedisCache {
	conn := &redis.Pool{
		MaxIdle:      conf.MaxIdle,
		MaxActive:    conf.MaxActive,
		IdleTimeout:  200 * time.Second,
		Dial:         getDial(conf.Host, conf.Password, conf.Port),
		TestOnBorrow: getTestOnBorrow(),
	}

	return &RedisCache{
		connection: conn,
		config:     conf,
	}
}

func (rc RedisCache) Exists(key string) (bool, error) {
	conn := rc.connection.Get()
	defer func() {
		err := conn.Close()
		if err != nil {
			panic(err)
		}
	}()

	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (rc RedisCache) Set(key string, data interface{}, time int) error {
	conn := rc.connection.Get()
	defer func() {
		err := conn.Close()
		if err != nil {
			panic(err)
		}
	}()

	value, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = conn.Do("SET", key, value)
	if err != nil {
		return err
	}

	_, err = conn.Do("EXPIRE", key, time)
	if err != nil {
		return err
	}

	return nil
}

func (rc RedisCache) Get(key string) ([]byte, error) {
	conn := rc.connection.Get()
	defer func() {
		err := conn.Close()
		if err != nil {
			panic(err)
		}
	}()

	reply, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return nil, err
	}

	return reply, nil
}

func (rc RedisCache) Delete(key string) (bool, error) {
	conn := rc.connection.Get()
	defer func() {
		err := conn.Close()
		if err != nil {
			panic(err)
		}
	}()

	return redis.Bool(conn.Do("DEL", key))
}

func (rc RedisCache) DeleteLike(key string) error {
	conn := rc.connection.Get()
	defer func() {
		err := conn.Close()
		if err != nil {
			panic(err)
		}
	}()

	keys, err := redis.Strings(conn.Do("KEYS", "*"+key+"*"))
	if err != nil {
		return err
	}

	for _, k := range keys {
		_, err = rc.Delete(k)
		if err != nil {
			return err
		}
	}

	return nil
}

func getTestOnBorrow() func(conn redis.Conn, t time.Time) error {
	return func(c redis.Conn, t time.Time) error {
		_, err := c.Do("PING")
		return err
	}
}

func getDial(host, password string, port int) func() (redis.Conn, error) {
	return func() (redis.Conn, error) {
		c, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
		if err != nil {
			return nil, err
		}

		if password != "" {
			_, err = c.Do("AUTH", password)

			if err != nil {
				err = c.Close()

				return nil, err
			}
		}
		return c, err
	}
}
