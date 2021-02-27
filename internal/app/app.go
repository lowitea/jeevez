package app

import (
	"fmt"
	"github.com/allegro/bigcache"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/kelseyhightower/envconfig"
	"github.com/lowitea/jeevez/internal/config"
	"github.com/lowitea/jeevez/internal/handlers"
	"github.com/lowitea/jeevez/internal/models"
	"github.com/lowitea/jeevez/internal/scheduler"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
	if err := db.AutoMigrate(&models.CurrencyRate{}); err != nil {
		log.Printf("migrating error: %s", err)
		os.Exit(1)
	}

	// запуск фоновых задач
	go scheduler.Run(bot, db, &cfg)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 1
	updates, _ := bot.GetUpdatesChan(u)

	// запуск обработки сообщений
	for update := range updates {
		go handlers.BaseCommandHandler(update, bot, &cfg)
		go handlers.CurrencyConverterHandler(update, bot, cache)
	}
}
