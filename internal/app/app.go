package app

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/lowitea/jeevez/internal/config"
	"github.com/lowitea/jeevez/internal/handlers"
	"github.com/lowitea/jeevez/internal/scheduler"
	"github.com/lowitea/jeevez/internal/tools"
	"log"
)

// Run функция запускающая бот
func Run() {
	// инициализируем конфиг
	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatalf("env parse error %s", err)
	}

	bot, err := tgbotapi.NewBotAPI(cfg.Telegram.Token)
	if err != nil {
		log.Fatalf("error connect to telegram %s", err)
	}
	log.Printf("Bot version: %s", cfg.App.Version)
	log.Printf("Authorized on account %s", bot.Self.UserName)

	// инициализируем кеш (пока не нужен)
	//cache, _ := bigcache.NewBigCache(bigcache.DefaultConfig(12 * time.Hour))

	// инициализация базы
	db, err := tools.InitDB(cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.DBName)
	if err != nil {
		log.Fatalf("error init database %s", err)
	}

	// запуск фоновых задач
	go scheduler.Run(bot, db)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 1
	updates, _ := bot.GetUpdatesChan(u)

	// отправляем инфу о запуске
	msg := tgbotapi.NewMessage(
		cfg.Telegram.Admin,
		"🤵🏻 Я обновился! :)\nМоя новая версия: "+cfg.App.Version,
	)
	_, _ = bot.Send(msg)

	// запуск обработки сообщений
	for update := range updates {
		// пропускаем, если сообщения нет
		if update.Message == nil {
			continue
		}
		go handlers.StartHandler(update, bot, db)
		go handlers.BaseCommandHandler(update, bot)
		go handlers.CurrencyConverterHandler(update, bot, db)
		go handlers.SwitchHandler(update, bot)
		go handlers.SubscriptionsHandler(update, bot, db)
		go handlers.DecorateTextHandler(update, bot)
	}
}
