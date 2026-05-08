package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	redis_client "github.com/ThuraMinThein/gateway/pkg/redis"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func SetCacheData(key string, value interface{}, expiration time.Duration) error {
	client := redis_client.GetRedisClient()
	if client == nil {
		return fmt.Errorf("redis client not initialized")
	}
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	err = client.Set(ctx, key, data, expiration).Err()
	return err
}

func GetCacheData(key string, result interface{}) error {
	client := redis_client.GetRedisClient()
	if client == nil {
		return fmt.Errorf("redis client not initialized")
	}
	val, err := client.Get(ctx, key).Result()
	if err == redis.Nil {
		return fmt.Errorf("cache miss")
	} else if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(val), result)
	return err
}

func InvalidateCache(key string) error {
	client := redis_client.GetRedisClient()
	if client == nil {
		return fmt.Errorf("redis client not initialized")
	}
	_, err := client.Del(ctx, key).Result()
	return err
}
