package handlers

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/lowitea/jeevez/internal/config"
	"math/rand"
	"strconv"
)

// cmdVersion вывод версии бота
func cmdVersion(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, config.Cfg.App.Version)
	_, _ = bot.Send(msg)
}

// cmdHelp вывод справки по боту
func cmdHelp(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	msgText := "" +
		"🤵🏻 Список команд:\n\n" +
		"/help - Показать этот список\n" +
		"/version - Показать текущую версию" +
		"\n\n" +
		"А ещё я слежу за текущим курсом и могу подсказать сколько стоит " +
		"нынче рубль :)\n" +
		"Напиши просто: `2000 долларов в рубли` или `1000 рублей в доллары`" +
		"\n\n" +
		"/currency_rate USD_RUB - Можно и просто курс текущий доллара к " +
		"рублю посмотреть например\n" +
		"/currency_rate - А так команда без параметров покажет все " +
		"доступные валютные пары)" +
		"\n\n" +
		"Также я могу сообщать вам полезную информацию, в удобное для вас время," +
		"только попросите:\n" +
		"/subscriptions - так я выведу список всех тем, о которых могу рассказать.\n" +
		"/subscribe covid19-moscow 11:00 - а так, можно задать интересующую Вас" +
		"тему, и время, в которое я буду Вам писать :)\n" +
		"/unsubscribe - так Вы сможете отменить подписку\n" +
		"/subscription covid19-moscow - с помощью этой команды можно получить сегодняшнею информацию, " +
		"без подписки на рассылку\n\n" +
		"/roll - а так я скажу Вам случайное число от 0 до 100\n" +
		"/switch L;bdbc - cfvsq kexibq ,jn!) - могу и текст Ваш починить)"
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
	_, _ = bot.Send(msg)
}

// cmdRoll выкидывает случайное число от 0 до 100
func cmdRoll(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	num := rand.Intn(101)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, strconv.Itoa(num))
	msg.ReplyToMessageID = update.Message.MessageID
	_, _ = bot.Send(msg)
}

// BaseCommandHandler базовый обработчик для выполнения команд. получает сообщение
// и если это сообщение - известная команда - вызывает нужную функцию
func BaseCommandHandler(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	// выходим сразу, если сообщения нет
	if update.Message == nil {
		return
	}

	cmdFuncMap := map[string]func(update tgbotapi.Update, bot *tgbotapi.BotAPI){
		"/version": cmdVersion,
		"/help":    cmdHelp,
		"/roll":    cmdRoll,
	}

	if cmdFunc, ok := cmdFuncMap[update.Message.Text]; ok {
		cmdFunc(update, bot)
		return
	}
}
