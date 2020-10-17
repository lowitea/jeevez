package handlers

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

func BaseHandler(update tgbotapi.Update, bot *tgbotapi.BotAPI) {

	if update.Message == nil { // ignore any non-Message Updates
		return
	}

	log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
	msg.ReplyToMessageID = update.Message.MessageID

	_, _ = bot.Send(msg)
}
