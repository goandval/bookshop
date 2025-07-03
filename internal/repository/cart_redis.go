package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type CartRedis struct {
	rdb *redis.Client
}

func NewCartRedis(rdb *redis.Client) *CartRedis {
	return &CartRedis{rdb: rdb}
}

func (r *CartRedis) SetReservation(ctx context.Context, userID string, bookID int, ttlSeconds int) error {
	key := fmt.Sprintf("reserve:book:%d", bookID)
	return r.rdb.Set(ctx, key, userID, time.Duration(ttlSeconds)*time.Second).Err()
}

func (r *CartRedis) RemoveReservation(ctx context.Context, userID string, bookID int) error {
	key := fmt.Sprintf("reserve:book:%d", bookID)
	return r.rdb.Del(ctx, key).Err()
}

func (r *CartRedis) IsReserved(ctx context.Context, bookID int) (bool, error) {
	key := fmt.Sprintf("reserve:book:%d", bookID)
	val, err := r.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return val != "", nil
}
