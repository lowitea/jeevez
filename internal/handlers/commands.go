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
		"/version - Показать текущую версию" +
		"\n\n" +
		"А ещё я слежу за текущим курсом и могу подсказать сколько стоит " +
		"нынче рубль :)\n" +
		"Напиши просто: `2000 долларов в рубли` или `1000 рублей в доллары`" +
		"\n\n" +
		"/̶c̶u̶r̶r̶e̶n̶c̶y̶R̶a̶t̶e̶ ̶-̶ ̶М̶о̶ж̶н̶о̶ ̶и̶ ̶п̶р̶о̶с̶т̶о̶ ̶к̶у̶р̶с̶ ̶т̶е̶к̶у̶щ̶и̶й̶ ̶п̶о̶с̶м̶о̶т̶р̶е̶т̶ь̶ (скоро)"
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
	_, _ = bot.Send(msg)
}

func cmdStart(update tgbotapi.Update, bot *tgbotapi.BotAPI, _ *config.Config) {
	msgText := "" +
		"Приветствую! Я Ваш личный бот помощник. 🤵🏻\n" +
		"Готов помогать всем, чем умею. Чтобы узнать, подробнее, " +
		"предлагаю использовать команду /help :)"
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
		"/start":   cmdStart,
	}

	if cmdFunc, ok := cmdFuncMap[update.Message.Text]; ok {
		cmdFunc(update, bot, cfg)
	}
}
