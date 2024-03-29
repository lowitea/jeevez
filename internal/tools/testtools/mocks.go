package testtools

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/mock"
)

// BotAPIMock мок объект для бота
type BotAPIMock struct {
	mock.Mock
	tgbotapi.BotAPI
}

// Send замоканный метод отправки сообщений в апи
func (s *BotAPIMock) Send(c tgbotapi.Chattable) (tgbotapi.Message, error) {
	args := s.Called(c)
	return args.Get(0).(tgbotapi.Message), args.Error(1) //nolint:forcetypeassert
}

// NewBotAPIMock возвращает настроенный мок для бота
func NewBotAPIMock(expMsg tgbotapi.Chattable) *BotAPIMock {
	botAPIMock := new(BotAPIMock)

	botAPIMock.On("Send", expMsg).Return(tgbotapi.Message{}, nil)

	return botAPIMock
}
