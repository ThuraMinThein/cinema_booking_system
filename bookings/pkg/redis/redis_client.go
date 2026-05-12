package redis_client

import (
	"context"
	"log"

	"github.com/ThuraMinThein/bookings/config"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var rdb *redis.Client

func InitRedis() error {

	rdb = redis.NewClient(&redis.Options{
		Addr:     config.Config.RedisURL,
		Username: config.Config.RedisUsername,
		Password: config.Config.RedisPassword,
		DB:       config.Config.RedisDB,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return err
	} else {
		log.Println("Redis connected successfully")
	}
	return nil
}

func GetRedisClient() *redis.Client {
	return rdb
}
