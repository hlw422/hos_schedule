package config

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func (c *RedisConfig) Connect() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", c.Host, c.Port),
		Password: c.Password,
		DB:       c.DB,
	})
}

func (c *RedisConfig) Ping(ctx context.Context, client *redis.Client) error {
	return client.Ping(ctx).Err()
}
