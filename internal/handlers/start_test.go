package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/lowitea/jeevez/internal/models"
	"github.com/lowitea/jeevez/internal/tools/testtools"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestStartHandler проверяет обработчки команды /start
func TestStartHandler(t *testing.T) {
	db := testtools.InitTestDB()
	db.Exec("DELETE FROM chats")

	successMsg := "Приветствую! Я Ваш личный бот помощник. 🤵🏻\n" +
		"Готов помогать всем, чем умею. Чтобы узнать, подробнее, " +
		"предлагаю использовать команду /help :)"

	// проверяем невалидную команду
	update := testtools.NewUpdate("/no_start")
	botAPIMock := testtools.NewBotAPIMock(tgbotapi.MessageConfig{})
	StartHandler(update, botAPIMock, db)
	botAPIMock.AssertNotCalled(t, "Send")

	// проверяем создание чата в базе
	update = testtools.NewUpdate("/start")
	update.Message.Chat.ID = 42
	expMsg := tgbotapi.NewMessage(update.Message.Chat.ID, successMsg)
	botAPIMock = testtools.NewBotAPIMock(expMsg)
	StartHandler(update, botAPIMock, db)
	botAPIMock.AssertExpectations(t)

	var chat models.Chat
	db.Last(&chat)

	assert.Equal(t, update.Message.Chat.ID, chat.TgID)

	// проверяем что ничего не падает при повторном создании чата с существующим id
	update = testtools.NewUpdate("/start")
	db.Create(&models.Chat{TgID: 1})
	expMsg = tgbotapi.NewMessage(update.Message.Chat.ID, successMsg)
	botAPIMock = testtools.NewBotAPIMock(expMsg)
	StartHandler(update, botAPIMock, db)
	botAPIMock.AssertExpectations(t)
}

// TestStartHandlerDBError проверяет работу хендлера при ошибке от базы данных
func TestStartHandlerDBError(t *testing.T) {
	db := testtools.InitTestDB()
	db.Exec("DROP TABLE chat_subscriptions")
	db.Exec("DROP TABLE chats")

	update := testtools.NewUpdate("/start")
	expMsg := tgbotapi.NewMessage(
		update.Message.Chat.ID,
		"Приветствую! Я Ваш личный бот помощник. 🤵🏻\n"+
			"К сожалению, не получилось Вас зарегистрировать, "+
			"попробуйте пожалуйста позже, с помощью команды /start ):",
	)
	botAPIMock := testtools.NewBotAPIMock(expMsg)
	StartHandler(update, botAPIMock, db)
	botAPIMock.AssertExpectations(t)
}
