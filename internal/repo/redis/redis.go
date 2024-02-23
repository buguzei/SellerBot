package redis

import (
	"bot/internal/config"
	log2 "bot/internal/log"
	"fmt"
	redis2 "github.com/redis/go-redis/v9"
)

type Redis struct {
	Client *redis2.Client
	logger log2.Logger
}

// NewRedis is a constructor for Redis
func NewRedis(cfg config.RedisConf, l log2.Logger) Redis {
	options := &redis2.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: "",
		DB:       0,
	}

	client := redis2.NewClient(options)

	l = l.Named("redis")

	return Redis{Client: client, logger: l}
}
