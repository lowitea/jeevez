package testFabrics

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

func NewUpdate(text string) tgbotapi.Update {
	chat := tgbotapi.Chat{ID: 1}
	msg := tgbotapi.Message{Chat: &chat, Text: text}
	return tgbotapi.Update{Message: &msg}
}
