package app

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/lowitea/jeevez/internal/config"
	"github.com/lowitea/jeevez/internal/handlers"
	"github.com/lowitea/jeevez/internal/models"
	"github.com/lowitea/jeevez/internal/scheduler"
	"github.com/lowitea/jeevez/internal/scheduler/subscriptions"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

// Run функция запускающая бот
func Run() {
	// инициализируем конфиг
	cfg, err := config.InitConfig()
	if err != nil {
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

	// инициализируем кеш (пока не нужен)
	//cache, _ := bigcache.NewBigCache(bigcache.DefaultConfig(12 * time.Hour))

	// инициализация базы
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d",
		cfg.DB.Host, cfg.DB.User, cfg.DB.Password, cfg.DB.DBName, cfg.DB.Port,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("failed to connect database: %s", err)
		os.Exit(1)
	}

	// миграция моделей
	if err := models.MigrateAll(db); err != nil {
		log.Printf("migrating error: %s", err)
		os.Exit(1)
	}

	// инициализация вариантов подписок
	if err := subscriptions.InitSubscriptions(db); err != nil {
		log.Printf("subscriptions init error: %s", err)
		os.Exit(1)
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
		go handlers.StartHandler(update, bot, db)
		go handlers.BaseCommandHandler(update, bot)
		go handlers.CurrencyConverterHandler(update, bot, db)
		go handlers.SwitchHandler(update, bot)
		go handlers.SubscriptionsHandler(update, bot, db)
		go handlers.DecorateTextHandler(update, bot)
	}
}
