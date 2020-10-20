package app

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/kelseyhightower/envconfig"
	"github.com/lowitea/jeevez/internal/handlers"
	"github.com/lowitea/jeevez/internal/scheduler"
	"log"
	"os"
)

func Run() {
	// инициализируем конфиг
	cfg := Config{}
	err := envconfig.Process("jeevez", &cfg)
	if err != nil {
		log.Printf("env parse error %s", err)
		os.Exit(1)
	}

	bot, err := tgbotapi.NewBotAPI(cfg.Telegram.Token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	// запуск фоновых задач
	go scheduler.Run(bot)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 1
	updates, _ := bot.GetUpdatesChan(u)

	// запуск обработки сообщений
	for update := range updates {
		go handlers.BaseHandler(update, bot)
		go handlers.VersionHandler(update, bot, cfg.App.Version)
	}
}
