package main

import (
	"bot/internal/config"
	"bot/internal/delivery/bot/telegram"
	"bot/internal/log"
	"bot/internal/repo"
	"bot/internal/repo/postgres"
	redis2 "bot/internal/repo/redis"
	"bot/internal/service"
)

func main() {
	var logger log.Logger = log.NewLogrus("debug")
	logger.Named("main")

	cfg, err := config.InitConfig()
	if err != nil {
		logger.Fatal("Initialization of config failed", log.Fields{
			"error": err,
		})
	}

	pg := postgres.NewPostgres(cfg.DB, logger)
	defer func() {
		err = pg.DB.Close()
		if err != nil {
			logger.Fatal("Closing Postgres failed", log.Fields{
				"error": err,
			})
		}
	}()

	redis := redis2.NewRedis(logger)
	defer func() {
		err = redis.Client.Close()
		if err != nil {
			logger.Fatal("Closing Redis failed", log.Fields{
				"error": err,
			})
		}
	}()

	repos := repo.Repo{
		UserRepo:  pg,
		OrderRepo: pg,
		CartRepo:  redis,
	}

	svc := service.NewService(repos, logger)

	bot := telegram.NewTGBot(svc, cfg.Bot, logger)
	bot.Run()
}
