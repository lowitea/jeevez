package app

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/lowitea/jeevez/internal/config"
	"github.com/lowitea/jeevez/internal/handlers"
	"github.com/lowitea/jeevez/internal/scheduler"
	"github.com/lowitea/jeevez/internal/structs"
	"github.com/lowitea/jeevez/internal/tools"
	"gorm.io/gorm"
	"log"
)

// processUpdate обрабатывает полученный апдейт
func processUpdate(update tgbotapi.Update, bot structs.Bot, db *gorm.DB) {
	// пропускаем, если сообщения нет
	if update.Message == nil {
		return
	}
	go handlers.StartHandler(update, bot, db)
	go handlers.BaseCommandHandler(update, bot)
	go handlers.CurrencyConverterHandler(update, bot, db)
	go handlers.SwitchHandler(update, bot)
	go handlers.SubscriptionsHandler(update, bot, db)
	go handlers.DecorateTextHandler(update, bot)
}

func initApp(
	initBotFunc func(token string) (*tgbotapi.BotAPI, error),
) (*tgbotapi.UpdatesChannel, structs.Bot, *gorm.DB, *config.Config, error) {
	// инициализируем конфиг
	cfg, err := config.InitConfig()
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("env parse error: %s", err)
	}

	bot, err := initBotFunc(cfg.Telegram.Token)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("error connect to telegram: %s", err)
	}
	log.Printf("Bot version: %s", cfg.App.Version)
	log.Printf("Authorized on account %s", bot.Self.UserName)

	// инициализация базы
	db, err := tools.InitDB(cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.DBName)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("error init database: %s", err)
	}

	// запуск фоновых задач
	go scheduler.Run(bot, db)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 1
	updates, _ := bot.GetUpdatesChan(u)

	return &updates, bot, db, cfg, nil
}

// releaseNotify отправляет сообщение админу о деплое
func releaseNotify(bot structs.Bot, adminID int64, version string) {
	msg := tgbotapi.NewMessage(
		adminID,
		fmt.Sprintf("🤵🏻 Я обновился! :)\nМоя новая версия: %s", version),
	)
	_, _ = bot.Send(msg)
}

// Run функция запускающая бот
func Run() {
	updates, bot, db, cfg, err := initApp(tgbotapi.NewBotAPI)
	if err != nil {
		log.Fatalf("error init app %s", err)
	}

	releaseNotify(bot, cfg.Telegram.Admin, cfg.App.Version)

	// запуск обработки сообщений
	for update := range *updates {
		processUpdate(update, bot, db)
	}
}
