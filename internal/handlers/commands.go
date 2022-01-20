package handlers

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/lowitea/jeevez/internal/config"
	"github.com/lowitea/jeevez/internal/structs"
	"math/rand"
	"strconv"
)

// cmdVersion вывод версии бота
func cmdVersion(update tgbotapi.Update, bot structs.Bot) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, config.Cfg.App.Version)
	_, _ = bot.Send(msg)
}

// cmdHelp вывод справки по боту
func cmdHelp(update tgbotapi.Update, bot structs.Bot) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, HelpText)
	_, _ = bot.Send(msg)
}

// cmdRoll выкидывает случайное число от 0 до 100
func cmdRoll(update tgbotapi.Update, bot structs.Bot) {
	num := rand.Intn(101) //nolint:gosec
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, strconv.Itoa(num))
	msg.ReplyToMessageID = update.Message.MessageID
	_, _ = bot.Send(msg)
}

func me(update tgbotapi.Update, bot structs.Bot) {
	msgText := fmt.Sprintf("Твой id: `%d`\nID чата: `%d`", update.Message.From.ID, update.Message.Chat.ID)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
	msg.ReplyToMessageID = update.Message.MessageID
	msg.ParseMode = "MARKDOWN"
	_, _ = bot.Send(msg)
}

var cmdFuncMap = map[string]func(update tgbotapi.Update, bot structs.Bot){
	"/version": cmdVersion,
	"/help":    cmdHelp,
	"/roll":    cmdRoll,
	"/me":      me,
}

// BaseCommandHandler базовый обработчик для выполнения команд. получает сообщение
// и если это сообщение - известная команда - вызывает нужную функцию
func BaseCommandHandler(update tgbotapi.Update, bot structs.Bot) {
	if cmdFunc, ok := cmdFuncMap[update.Message.Text]; ok {
		cmdFunc(update, bot)
		return
	}
}
