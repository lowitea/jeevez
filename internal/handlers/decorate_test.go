package handlers

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/lowitea/jeevez/internal/tools/testTools"
	"testing"
)

// TestDecorateTextHandler проверяет обработчик команд для обработки текста
func TestDecorateTextHandler(t *testing.T) {
	cases := [...]struct {
		Cmd     string
		MsgText string
	}{
		{"strth test", "t̶e̶s̶t̶"},
		{"strth русский текст", "р̶у̶с̶с̶к̶и̶й̶ ̶т̶е̶к̶с̶т̶"},
		{"reverse test", "tset"},
		{"reverse русский текст", "тскет йикссур"},
		{"invert test", "ʇsǝʇ"},
		{"invert русский текст", "ɯɔʞǝɯ ņиʞɔɔʎd"},
		{"zalgo test", "t̘̤̠͔̦̝͙̼̬̼̠̯̟̟̠͉̽͋́ͪ̆̍͢͡ͅę͙̗͈̣̖̻̠̯̹͕̦͉̝͖̼͇̗ͧ̽̓̆ͧ̽ͥ̋̊ͪͮ̊ͣ͛̒̀͘ş̹͖̱̹͓̮̪̗̰͓̙̊̋̄͘͞͡t̹͎̙̝͚̻͍͕͙̟̥̥̝͉͎̯͓͍ͦ̾̐͡"},
		{"zalgo русский текст", "р͊ͧ̂ͮ҉̣̼͉̕̕у̴̷̩̭̣̱̼̦͖̖͍̽̄́ͥ͑ͧ̒͋ͤ̚͜͟с̡̬̦̣͉͓̱̣̖̹͉̺̟̥̝̘ͣ̎ͪ̄͒͊̽̒̿̌͐ͬ͒ͭ͋̈́͑ͥ̕͠͝с̶̣̘̟̆ͯ̈́̍̆̓ͬ̽̏ͦк̲̪͕̅ͨ͋̆͌̐̐̄͜и̧͙͖̖̣̾ͦ̉̃͑̌̓̇͑̃̓̿͂͑̎й̲͉̫͈͉̿ͥ͗̓͌͂ͪ̈ͪ̍́͝͝͝ ̧̛̻̘̣͎̤͍̗̬͍̹̪ͮͯ̀ͧ̀ͭ̆͆̄̃̌̀т̮̫̣̦̼̹̠͖̞͙̩͖̦̻̘͈̯͒͒̒̅̀ͤ́̀̆ͧͦ͊͂͒͗̾͐̀͟е̸͍̪̖̖̗̝͖͈̼͖̬̱͔̒͗̇̄͠ќ̡̛̪͚̟̭̠͓̦̫̼̱̦̩ͭ̔͒̽̓͒ͤͬ͐̀с̨̘̝̫̞̣̺̱̜̣̟̜͙ͣ̔ͦͧ̈̿̂ͭ͜͟т̧̛̩̩͕̮̳̪̮̣̟̱̝̞͔ͪͩ͂ͯͩ̽̽̎͘͠"},
		{"no_func test", "К сожалению, такому меня не учили ):"},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("cmd=%s", c.Cmd), func(t *testing.T) {
			update := testTools.NewUpdate("/decorate " + c.Cmd)
			expMsg := tgbotapi.NewMessage(update.Message.Chat.ID, c.MsgText)
			botAPIMock := testTools.NewBotAPIMock(expMsg)

			DecorateTextHandler(update, botAPIMock)

			botAPIMock.AssertExpectations(t)
		})
	}
}

// TestDecorateTextHandlerBadCmd проверяет обработку невалидной команды
func TestDecorateTextHandlerBadCmd(t *testing.T) {
	update := testTools.NewUpdate("невалидная команда")
	botAPIMock := testTools.NewBotAPIMock(tgbotapi.MessageConfig{})

	DecorateTextHandler(update, botAPIMock)

	botAPIMock.AssertNotCalled(t, "Send")
}

// TestDecorateTextHandlerAllCommands проверяет вывод списка комманд
func TestDecorateTextHandlerAllCommands(t *testing.T) {
	update := testTools.NewUpdate("/decorate")
	expMsg := tgbotapi.NewMessage(
		update.Message.Chat.ID,
		"Чтобы попросить меня декорировать текст, отправьте команду в формате:\n"+
			"/decorate название_преобразования Ваш текст\n"+
			"Например, так:\n"+
			"/decorate strth Текст, который нужно зачеркнуть\n\n"+
			"Все варианты декорирования:\n\n"+
			"  <b>invert</b> - Перевёрнутый текст\n"+
			"  <b>reverse</b> - Текст в обратную сторону\n"+
			"  <b>strth</b> - Зачёркнутый текст\n"+
			"  <b>zalgo</b> - Зальгофикация текста",
	)
	expMsg.ParseMode = "HTML"
	botAPIMock := testTools.NewBotAPIMock(expMsg)

	DecorateTextHandler(update, botAPIMock)

	botAPIMock.AssertExpectations(t)
}
