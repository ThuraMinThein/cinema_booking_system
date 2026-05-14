package memory_storage

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ThuraMinThein/bookings/internal/model"
	redis_client "github.com/ThuraMinThein/bookings/pkg/redis"
	"github.com/ThuraMinThein/common/api"
	"github.com/redis/go-redis/v9"
)

func SetData(key string, value interface{}, expiration time.Duration) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

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

func SaveBooking(request *api.HoldBookingRequest, booking interface{}) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client := redis_client.GetRedisClient()

	key := fmt.Sprintf("booking:%d:%d", request.MovieId, request.SeatId)
	indexKey := fmt.Sprintf("movie_bookings:%d", request.MovieId)

	data, err := json.Marshal(booking)
	if err != nil {
		return err
	}

	pipe := client.TxPipeline()

	pipe.Set(ctx, key, data, 2*time.Minute)

	pipe.SAdd(ctx, indexKey, key)

	pipe.Expire(ctx, indexKey, 2*time.Minute)

	_, err = pipe.Exec(ctx)

	return err
}

func GetData(key string, result interface{}) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client := redis_client.GetRedisClient()
	if client == nil {
		return fmt.Errorf("redis client not initialized")
	}
	val, err := client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil
	} else if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(val), result)
	return err
}

func GetMovieBookings(movieId int64) ([]model.Booking, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client := redis_client.GetRedisClient()

	indexKey := fmt.Sprintf("movie_bookings:%d", movieId)

	keys, err := client.SMembers(ctx, indexKey).Result()
	if err != nil {
		return nil, err
	}

	if len(keys) == 0 {
		return []model.Booking{}, nil
	}

	values, err := client.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}

	bookings := make([]model.Booking, 0)

	for _, val := range values {
		if val == nil {
			continue
		}

		var booking model.Booking

		strVal, ok := val.(string)
		if !ok {
			continue
		}

		if err := json.Unmarshal([]byte(strVal), &booking); err != nil {
			continue
		}

		bookings = append(bookings, booking)
	}

	return bookings, nil
}

func InvalidateStorage(key string) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client := redis_client.GetRedisClient()
	if client == nil {
		return fmt.Errorf("redis client not initialized")
	}
	_, err := client.Del(ctx, key).Result()
	return err
}
