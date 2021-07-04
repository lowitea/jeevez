package handlers

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/lowitea/jeevez/internal/config"
	"github.com/lowitea/jeevez/internal/tools/testTools"
	"testing"
)

// TestBaseCommandHandler проверяет обработчик базовых команд
func TestBaseCommandHandler(t *testing.T) {
	cases := [...]struct {
		Cmd     string
		MsgText string
	}{
		{"/version", config.Cfg.App.Version},
		{"/help", HelpText},
		{"/roll", "65"},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("cmd=%s", c.Cmd), func(t *testing.T) {
			update := testTools.NewUpdate(c.Cmd)
			expMsg := tgbotapi.NewMessage(update.Message.Chat.ID, c.MsgText)
			botAPIMock := testTools.NewBotAPIMock(expMsg)

			BaseCommandHandler(update, botAPIMock)

			botAPIMock.AssertExpectations(t)
		})
	}

	// при неизвестной команде не падаем и ничего не делаем
	t.Run("cmd=/unknown", func(t *testing.T) {
		update := testTools.NewUpdate("/unknown")
		botAPIMock := testTools.NewBotAPIMock(tgbotapi.MessageConfig{})

		BaseCommandHandler(update, botAPIMock)

		botAPIMock.AssertNotCalled(t, "Send")
	})
}
