package app

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/lowitea/jeevez/internal/tools/testTools"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

// TestInitApp смоук тест инициализации приложения
func TestInitApp(t *testing.T) {
	dbHost := os.Getenv("JEEVEZ_TEST_DB_HOST")

	_ = os.Setenv("JEEVEZ_TELEGRAM_TOKEN", "1")
	_ = os.Setenv("JEEVEZ_TELEGRAM_ADMIN", "1")
	_ = os.Setenv("JEEVEZ_DB_USER", "test")
	_ = os.Setenv("JEEVEZ_DB_PASSWORD", "test")
	_ = os.Setenv("JEEVEZ_DB_HOST", dbHost)
	_ = os.Setenv("JEEVEZ_CURRENCYAPI_TOKEN", "1")

	initBotFunc := func(_ string) (*tgbotapi.BotAPI, error) { return &tgbotapi.BotAPI{}, nil }

	assert.NotPanics(t, func() { _, _, _, _, _ = initApp(initBotFunc) })
}

// TestReleaseNotify проверяет функцию отправки сообщения о релизе
func TestReleaseNotify(t *testing.T) {
	var adminID int64 = 666
	expMsg := tgbotapi.NewMessage(adminID, "🤵🏻 Я обновился! :)\nМоя новая версия: 6.6.6")
	botAPIMock := testTools.NewBotAPIMock(expMsg)
	releaseNotify(botAPIMock, adminID, "6.6.6")
	botAPIMock.AssertExpectations(t)
}

// TestProcessUpdate смоук тест общего запуска хендлеров
func TestProcessUpdate(t *testing.T) {
	db := testTools.InitTestDB()
	update := testTools.NewUpdate("no_command")
	botAPIMock := testTools.NewBotAPIMock(tgbotapi.MessageConfig{})
	assert.NotPanics(t, func() { processUpdate(update, botAPIMock, db) })
	botAPIMock.AssertNotCalled(t, "Send")

	assert.NotPanics(t, func() { processUpdate(tgbotapi.Update{}, botAPIMock, db) })
	botAPIMock.AssertNotCalled(t, "Send")
}
