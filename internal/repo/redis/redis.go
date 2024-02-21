package redis

import (
	log2 "bot/internal/log"
	redis2 "github.com/redis/go-redis/v9"
)

type Redis struct {
	Client *redis2.Client
	logger log2.Logger
}

// NewRedis is a constructor for Redis
func NewRedis(l log2.Logger) Redis {
	options := &redis2.Options{
		Addr:     "redis:6380",
		Password: "",
		DB:       0,
	}

	client := redis2.NewClient(options)

	l = l.Named("redis")

	return Redis{Client: client, logger: l}
}
