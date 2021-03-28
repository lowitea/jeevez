package structs

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

type Bot interface {
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
}
