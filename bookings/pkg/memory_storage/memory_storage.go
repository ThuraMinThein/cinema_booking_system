package memory_storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/ThuraMinThein/bookings/internal/model"
	redis_client "github.com/ThuraMinThein/bookings/pkg/redis"
	"github.com/ThuraMinThein/common/api"
	"github.com/redis/go-redis/v9"
)

func SetData(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
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

func SaveBooking(ctx context.Context, request *api.HoldBookingRequest, booking interface{}) error {

	client := redis_client.GetRedisClient()

	key := fmt.Sprintf("booking:%d:%d", request.MovieId, request.SeatId)
	indexKey := fmt.Sprintf("movie_bookings:%d", request.MovieId)

	data, err := json.Marshal(booking)
	if err != nil {
		return err
	}

	pipe := client.TxPipeline()

	setCmd := pipe.SetNX(ctx, key, data, 2*time.Minute)

	pipe.SAdd(ctx, indexKey, key)

	pipe.Expire(ctx, indexKey, 2*time.Minute)

	_, err = pipe.Exec(ctx)
	if err != nil {
		return err
	}

	ok, err := setCmd.Result()
	if err != nil {
		return err
	}

	if !ok {
		return errors.New("seat already held")
	}

	return nil
}

func GetData(ctx context.Context, key string, result interface{}) error {

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

func GetMovieBookings(ctx context.Context, movieId int64) ([]model.Booking, error) {

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

func InvalidateStorage(ctx context.Context, key string) error {

	client := redis_client.GetRedisClient()
	if client == nil {
		return fmt.Errorf("redis client not initialized")
	}
	_, err := client.Del(ctx, key).Result()
	return err
}
