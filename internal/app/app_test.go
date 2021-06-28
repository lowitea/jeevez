package app

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
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
