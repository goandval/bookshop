package integration

import (
	"context"
)

type KeycloakClient interface {
	ValidateToken(ctx context.Context, token string) (userID, email string, roles []string, err error)
}

type RedisCache interface {
	Get(key string) (string, error)
	Set(key string, value string, ttlSeconds int) error
	Del(key string) error
}

type KafkaProducer interface {
	PublishOrderPlaced(ctx context.Context, orderID int, userID string, bookIDs []int) error
}
