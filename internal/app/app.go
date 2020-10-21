package app

import (
	"github.com/allegro/bigcache"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/kelseyhightower/envconfig"
	"github.com/lowitea/jeevez/internal/config"
	"github.com/lowitea/jeevez/internal/handlers"
	"github.com/lowitea/jeevez/internal/scheduler"
	"log"
	"os"
	"time"
)

// Run функция запускающая бот
func Run() {
	// инициализируем конфиг
	cfg := config.Config{}

	if err := envconfig.Process("jeevez", &cfg); err != nil {
		log.Printf("env parse error %s", err)
		os.Exit(1)
	}

	bot, err := tgbotapi.NewBotAPI(cfg.Telegram.Token)
	if err != nil {
		log.Printf("error connect to telegram %s", err)
		os.Exit(1)
	}
	log.Printf("Bot version: %s", cfg.App.Version)
	log.Printf("Authorized on account %s", bot.Self.UserName)

	// отправляем инфу о запуске
	msg := tgbotapi.NewMessage(
		cfg.Telegram.Admin,
		"🤵🏻 Я обновился! :)\nМоя новая версия: "+cfg.App.Version,
	)
	_, _ = bot.Send(msg)

	// инициализируем кеш
	cache, _ := bigcache.NewBigCache(bigcache.DefaultConfig(12 * time.Hour))

	// запуск фоновых задач
	go scheduler.Run(bot)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 1
	updates, _ := bot.GetUpdatesChan(u)

	// запуск обработки сообщений
	for update := range updates {
		go handlers.BaseCommandHandler(update, bot, &cfg)
		go handlers.CurrencyConverterHandler(update, bot, cache)
	}
}
