package handlers

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/lowitea/jeevez/internal/config"
)

// cmdVersion вывод версии бота
func cmdVersion(update tgbotapi.Update, bot *tgbotapi.BotAPI, cfg *config.Config) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, cfg.App.Version)
	_, _ = bot.Send(msg)
}

// cmdHelp вывод справки по боту
func cmdHelp(update tgbotapi.Update, bot *tgbotapi.BotAPI, _ *config.Config) {
	msgText := "" +
		"🤵🏻 Список команд:\n\n" +
		"/help - Показать этот список\n" +
		"/version - Показать текущую версию"
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
	_, _ = bot.Send(msg)
}

// BaseCommandHandler базовый обработчик для выполнения команд. получает сообщение
// и если это сообщение - известная команда - вызывает нужную функцию
func BaseCommandHandler(update tgbotapi.Update, bot *tgbotapi.BotAPI, cfg *config.Config) {
	// выходим сразу, если сообщения нет
	if update.Message == nil {
		return
	}

	cmdFuncMap := map[string]func(update tgbotapi.Update, bot *tgbotapi.BotAPI, cfg *config.Config){
		"/version": cmdVersion,
		"/help":    cmdHelp,
	}

	if cmdFunc, ok := cmdFuncMap[update.Message.Text]; ok {
		cmdFunc(update, bot, cfg)
	}
}
