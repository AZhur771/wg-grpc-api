package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v9"
)

const redisTimeout time.Duration = 3 * time.Second

func New(ctx context.Context, redisAddr string, redisPassword string) (*redis.Client, error) {
	ctx, cancel := context.WithTimeout(ctx, redisTimeout)
	defer cancel()

	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,     // host:port of the redis server
		Password: redisPassword, // no password set
		DB:       0,             // use default DB
	})

	cmd := client.Ping(ctx)

	return client, cmd.Err()
}
