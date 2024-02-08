package redis

import redis2 "github.com/redis/go-redis/v9"

type Redis struct {
	Client *redis2.Client
}

func NewRedis() Redis {
	options := &redis2.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	}

	client := redis2.NewClient(options)

	return Redis{Client: client}
}
