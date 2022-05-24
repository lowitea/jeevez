package app

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/lowitea/jeevez/internal/config"
	"github.com/lowitea/jeevez/internal/handlers"
	"github.com/lowitea/jeevez/internal/scheduler"
	"github.com/lowitea/jeevez/internal/structs"
	"github.com/lowitea/jeevez/internal/tools"
	"gorm.io/gorm"
	"log"
	"strings"
)

type appDepContainer struct {
	updateConfig *tgbotapi.UpdateConfig
	bot          structs.Bot
	db           *gorm.DB
	cfg          *config.Config
}

// checkRecipient проверяет что сообщение адресовано именно этому боту
func checkRecipient(update *tgbotapi.Update, botName string) bool {
	args := strings.Split(update.Message.Text, "@")

	if update.Message.Chat.IsPrivate() {
		update.Message.Text = args[0]
		return true
	}

	if len(args) != 2 || args[1] != botName {
		return false
	}
	update.Message.Text = args[0]
	return true
}

// processUpdate обрабатывает полученный апдейт
func processUpdate(update tgbotapi.Update, dep *appDepContainer) {
	// пропускаем, если сообщения нет
	if update.Message != nil {
		// пропускаем, если сообщение отправлено не нам
		if !checkRecipient(&update, dep.cfg.Telegram.BotName) {
			return
		}

		go handlers.UpdateChatInfoHandler(update, dep.db)
		go handlers.StartHandler(update, dep.bot, dep.db)
		go handlers.BaseCommandHandler(update, dep.bot)
		go handlers.CurrencyConverterHandler(update, dep.bot, dep.db)
		go handlers.SwitchHandler(update, dep.bot)
		go handlers.SubscriptionsHandler(update, dep.bot, dep.db)
		go handlers.DecorateTextHandler(update, dep.bot)
		go handlers.YogaHandler(update, dep.bot)
		go handlers.WeatherHandler(&update, dep.bot)

		// команды для админского аккаунта
		if update.Message.From.ID == dep.cfg.Telegram.Admin {
			go handlers.ChatListHandler(update, dep.bot, dep.db)
			go handlers.DeleteChatHandler(update, dep.bot, dep.db)
		}
	} else if update.CallbackQuery != nil {
		go handlers.YogaCallbackHandler(update, dep.bot)
	}
}

func initApp(
	initBotFunc func(token string) (*tgbotapi.BotAPI, error),
	initCfgFunc func() (*config.Config, error),
) (*appDepContainer, error) {
	// инициализируем конфиг
	cfg, err := initCfgFunc()
	if err != nil {
		return nil, fmt.Errorf("env parse error: %w", err)
	}

	bot, err := initBotFunc(cfg.Telegram.Token)
	if err != nil {
		return nil, fmt.Errorf("error connect to telegram: %w", err)
	}
	log.Printf("Bot version: %s", cfg.App.Version)
	log.Printf("Authorized on account %s", bot.Self.UserName)

	// инициализация базы
	db, err := tools.InitDB(cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Name)
	if err != nil {
		return nil, fmt.Errorf("error init database: %w", err)
	}

	// запуск фоновых задач
	scheduler.Run(bot, db)

	updateCfg := tgbotapi.NewUpdate(0)
	updateCfg.Timeout = 1

	return &appDepContainer{updateConfig: &updateCfg, bot: bot, db: db, cfg: cfg}, nil
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
	depContainer, err := initApp(tgbotapi.NewBotAPI, config.InitConfig)
	if err != nil {
		log.Panicf("error init app %s\n", err)
	}

	updates := depContainer.bot.GetUpdatesChan(*depContainer.updateConfig)
	releaseNotify(depContainer.bot, depContainer.cfg.Telegram.Admin, depContainer.cfg.App.Version)

	// запуск обработки сообщений
	for update := range updates {
		processUpdate(update, depContainer)
	}
}
