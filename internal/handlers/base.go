package handlers

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

func VersionHandler(update tgbotapi.Update, bot *tgbotapi.BotAPI, version string) {
	if update.Message == nil { // ignore any non-Message Updates
		return
	}

	if update.Message.Text != "/version" {
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, version)
	_, _ = bot.Send(msg)
}

func BaseHandler(update tgbotapi.Update, bot *tgbotapi.BotAPI) {

	if update.Message == nil { // ignore any non-Message Updates
		return
	}

	log.Printf("[%s (%d)] %s", update.Message.From.UserName, update.Message.Chat.ID, update.Message.Text)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
	msg.ReplyToMessageID = update.Message.MessageID

	//_, _ = bot.Send(msg)
}
