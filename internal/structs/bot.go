package structs

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

type Bot interface {
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
	GetUpdatesChan(config tgbotapi.UpdateConfig) (tgbotapi.UpdatesChannel, error)
	AnswerCallbackQuery(config tgbotapi.CallbackConfig) (tgbotapi.APIResponse, error)
}
