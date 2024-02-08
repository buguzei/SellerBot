package telegram

import (
	"bot/internal/config"
	log2 "bot/internal/log"
	"bot/internal/service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TGBot struct {
	conf   config.BotConf
	bot    *tgbotapi.BotAPI
	svc    service.Service
	cache  map[int64]map[string]interface{}
	logger log2.Logger
}

func (tg TGBot) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	tg.bot.Debug = false

	updates := tg.bot.GetUpdatesChan(u)
	for update := range updates {
		tg.HandleUpdates(update)
	}
}

func (tg TGBot) HandleUpdates(update tgbotapi.Update) {
	if update.Message != nil {
		tg.HandleMessage(update)
	}

	if update.CallbackQuery != nil {
		tg.HandleCallback(update)
	}
}

func NewTGBot(svc service.Service, conf config.BotConf, l log2.Logger) TGBot {
	bot, err := tgbotapi.NewBotAPI(conf.Token)
	if err != nil {
		panic(err)
	}

	cache := make(map[int64]map[string]interface{})

	l.Named("bot")

	l.Info("New tg bot", log2.Fields{
		"name": bot.Self.UserName,
	})

	return TGBot{
		bot:    bot,
		svc:    svc,
		logger: l,
		cache:  cache,
		conf:   conf,
	}
}
