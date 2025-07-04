package integration

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCacheImpl struct {
	rdb *redis.Client
}

func NewRedisCache(rdb *redis.Client) *RedisCacheImpl {
	return &RedisCacheImpl{rdb: rdb}
}

func (r *RedisCacheImpl) Get(key string) (string, error) {
	ctx := context.Background()
	return r.rdb.Get(ctx, key).Result()
}

func (r *RedisCacheImpl) Set(key string, value string, ttlSeconds int) error {
	ctx := context.Background()
	return r.rdb.Set(ctx, key, value, time.Duration(ttlSeconds)*time.Second).Err()
}

func (r *RedisCacheImpl) Del(key string) error {
	ctx := context.Background()
	return r.rdb.Del(ctx, key).Err()
}
