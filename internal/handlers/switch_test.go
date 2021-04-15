package handlers

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/lowitea/jeevez/internal/tools/testTools"
	"testing"
)

// TestSwitchHandler проверяет обработчик команд исправления раскладки текста
func TestSwitchHandler(t *testing.T) {
	cases := [...]struct {
		Cmd     string
		MsgText string
	}{
		// не баг, а фича, так работает автоопределитель раскладки
		{"/switch тест", "ntcn"},
		{".ыцшеср test", "еуые"},

		// а это уже норм кейсы
		{"/switch ntcn", "тест"},
		{".ыцшеср еуые", "test"},

		// проверяем ситуацию когда смешаны раскладки
		{
			"/switch 'nj heccrbq NTRCN не с одним fyukbqcrbv ckjdjv",
			"это русский ТЕКСТ не с одним английским словом",
		},
		{
			".ыцшеср Ше шы ф утпдшыр еуче цшер one кгыышфт цщкв",
			"It is a english text with one russian word",
		},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("cmd=%s", c.Cmd), func(t *testing.T) {
			update := testTools.NewUpdate(c.Cmd)
			expMsg := tgbotapi.NewMessage(update.Message.Chat.ID, c.MsgText)
			botAPIMock := testTools.NewBotAPIMock(expMsg)

			SwitchHandler(update, botAPIMock)

			botAPIMock.AssertExpectations(t)
		})
	}
}

// TestSwitchHandlerInvalidCmd проверяет обработчик при невалидных командах
func TestSwitchHandlerInvalidCmd(t *testing.T) {
	update := testTools.NewUpdate("/no_switch")
	botAPIMock := testTools.NewBotAPIMock(tgbotapi.MessageConfig{})

	SwitchHandler(update, botAPIMock)

	botAPIMock.AssertNotCalled(t, "Send")

	// проверяем команду без параметров
	update = testTools.NewUpdate("/switch")
	expMsg := tgbotapi.NewMessage(
		update.Message.Chat.ID,
		"Пришлите текст для изменения его раскладки.\n"+
			"Пример команды:\n"+
			"/switch L;bdbc - cfvsq kexibq ,jn!)",
	)
	botAPIMock = testTools.NewBotAPIMock(expMsg)

	SwitchHandler(update, botAPIMock)

	botAPIMock.AssertExpectations(t)
}
