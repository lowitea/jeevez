package testtools

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func NewUpdate(text string) tgbotapi.Update {
	chat := tgbotapi.Chat{ID: 1}
	msg := tgbotapi.Message{Chat: &chat, Text: text, From: &tgbotapi.User{}}
	return tgbotapi.Update{Message: &msg}
}
