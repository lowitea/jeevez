package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/lowitea/jeevez/internal/models"
	"github.com/lowitea/jeevez/internal/tools/testTools"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestStartHandler проверяет обработчки команды /start
func TestStartHandler(t *testing.T) {
	db, _ := testTools.InitTestDB()

	successMsg := "Приветствую! Я Ваш личный бот помощник. 🤵🏻\n" +
		"Готов помогать всем, чем умею. Чтобы узнать, подробнее, " +
		"предлагаю использовать команду /help :)"

	// проверяем невалидную команду
	update := testTools.NewUpdate("/no_start")
	botAPIMock := testTools.NewBotAPIMock(tgbotapi.MessageConfig{})

	StartHandler(update, botAPIMock, db)

	botAPIMock.AssertNotCalled(t, "Send")

	// проверяем создание чата в базе
	update = testTools.NewUpdate("/start")
	update.Message.Chat.ID = 42
	expMsg := tgbotapi.NewMessage(update.Message.Chat.ID, successMsg)
	botAPIMock = testTools.NewBotAPIMock(expMsg)

	StartHandler(update, botAPIMock, db)

	botAPIMock.AssertExpectations(t)

	var chat models.Chat
	db.Last(&chat)

	assert.Equal(t, update.Message.Chat.ID, chat.TgID)

	// проверяем что ничего не падает при повторном создании чата с существующим id
	update = testTools.NewUpdate("/start")
	db.Create(&models.Chat{TgID: 1})
	expMsg = tgbotapi.NewMessage(update.Message.Chat.ID, successMsg)
	botAPIMock = testTools.NewBotAPIMock(expMsg)

	StartHandler(update, botAPIMock, db)

	botAPIMock.AssertExpectations(t)
}
