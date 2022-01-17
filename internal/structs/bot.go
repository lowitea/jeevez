package structs

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Bot interface {
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
	GetUpdatesChan(config tgbotapi.UpdateConfig) tgbotapi.UpdatesChannel
	Request(c tgbotapi.Chattable) (*tgbotapi.APIResponse, error)
}
