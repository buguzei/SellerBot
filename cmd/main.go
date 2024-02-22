package main

import (
	"bot/internal/config"
	"bot/internal/delivery/bot/telegram"
	"bot/internal/log"
	"bot/internal/repo"
	"bot/internal/repo/postgres"
	redis2 "bot/internal/repo/redis"
	"bot/internal/service"
	"github.com/pressly/goose"
)

func main() {
	// init Logger
	var logger log.Logger = log.NewLogrus("debug")
	logger.Named("main")

	// init config
	cfg, err := config.InitConfig()
	if err != nil {
		logger.Fatal("Initialization of config failed", log.Fields{
			"error": err,
		})
	}

	// init postgres
	pg := postgres.NewPostgres(cfg.DB, logger)
	defer func() {
		err = pg.DB.Close()
		if err != nil {
			logger.Fatal("main: Closing Postgres failed", log.Fields{
				"error": err,
			})
		}
	}()

	// make migrations
	err = goose.Up(pg.DB, "./db/")
	if err != nil {
		logger.Fatal("main: goose up failed", log.Fields{
			"error": err,
		})
	}

	// init redis
	redis := redis2.NewRedis(logger)
	defer func() {
		err = redis.Client.Close()
		if err != nil {
			logger.Fatal("main: Closing Redis failed", log.Fields{
				"error": err,
			})
		}
	}()

	// create repo
	repos := repo.Repo{
		UserRepo:  pg,
		OrderRepo: pg,
		CartRepo:  redis,
	}

	// create svc
	svc := service.NewService(repos, logger)

	// running bot
	bot := telegram.NewTGBot(svc, cfg.Bot, logger)
	bot.Run()
}
