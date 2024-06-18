package database

import (
	"sync"

	"github.com/fredyranthun/url-shortner/config"
	"github.com/redis/go-redis/v9"
)

var (
  once       sync.Once
  redisClient *redis.Client
)

func NewRedisClient(conf *config.Config) *redis.Client {
	once.Do(func() {
		redisClient = redis.NewClient(&redis.Options{
			Addr: conf.RdbAddr,
			Password: conf.RdbPassword,
			DB: 0,
		})	
	})

	return redisClient
}